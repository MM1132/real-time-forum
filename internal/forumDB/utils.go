package forumDB

import (
	"database/sql"
	"fmt"
	"forum/internal/utils"
	"log"
	"os"
	"regexp"
	"strings"
)

func makeStatementMap(db *sql.DB, sqlPath string) map[string]*sql.Stmt {
	statements := make(map[string]*sql.Stmt)

	file, err := os.ReadFile(sqlPath)
	utils.FatalErr(err)

	// Semicolon separates SQL statements
	stmtSlice := strings.SplitAfter(string(file), ";")
	stmtSlice = stmtSlice[:len(stmtSlice)-1] // Remove last

	r := regexp.MustCompile(`(?i)\s*--\s*Func:\s*(\S+)`)
	rs := regexp.MustCompile(`{(.+?)}`)
	rCount1 := regexp.MustCompile(`-- ?WITH COUNT`)
	rCount2 := regexp.MustCompile(`(?s)ORDER BY.*`)

	for i := 0; i < len(stmtSlice); i++ {
		stmt := stmtSlice[i]
		nameSlc := r.FindStringSubmatch(stmt)
		if nameSlc == nil {
			log.Panicf("Could not find a function name in %v:\n%v\n", sqlPath, stmt)
		}

		name := nameSlc[1]

		if rCount1.MatchString(stmt) {
			// Create a counting statement
			tempStmt := stmt
			tempStmt = r.ReplaceAllString(tempStmt, fmt.Sprintf(`-- Func: %vCount`, name[:strings.Index(name, "{")]))
			tempStmt = rCount1.ReplaceAllString(tempStmt, "")
			tempStmt = rCount2.ReplaceAllString(tempStmt, ";")
			tempStmt = "SELECT count(*) FROM (" + tempStmt[:len(tempStmt)-1] + ");"
			stmtSlice = append(stmtSlice, tempStmt)
		}

		if matches := rs.FindAllStringSubmatch(name, -1); matches == nil {
			statements[name], err = db.Prepare(stmt)
			if err != nil {
				log.Panic(fmt.Errorf("error preparing %v: %w", name, err))
			}
		} else {
			name = name[:strings.Index(name, "{")]

			var signatures [][]string
			for _, matchSlc := range matches {
				match := matchSlc[1]

				vars := strings.Split(match, ",")
				signatures = append(signatures, vars)
			}

			specialStatement(db, statements, name, stmt, signatures)
		}
	}

	return statements
}

// For every combination of signatures, make a new statement where the text between '--START' and '--END' has been replaced
// with a space separated combination of signature strings.
// The statement's name will be suffixed by the combination used.
func specialStatement(db *sql.DB, statements map[string]*sql.Stmt, stmtName, stmtString string, signatures [][]string) {
	reg := regexp.MustCompile(`(?s)-- ?START.*-- ?END`)

	variations := cartN(signatures...)

	for _, va := range variations {
		key := strings.Join(va, "-")
		value := strings.Join(va, " ")
		newStmtString := reg.ReplaceAllString(stmtString, value)

		var err error
		statements[stmtName+"-"+key], err = db.Prepare(newStmtString)
		if err != nil {
			log.Panic(fmt.Errorf("error preparing %v: %w", key, err))
		}
	}
}

// Arranges cartesian product
func cartN(a ...[]string) (c [][]string) {
	if len(a) == 0 {
		return [][]string{nil}
	}
	r := cartN(a[1:]...)
	for _, e := range a[0] {
		for _, p := range r {
			c = append(c, append([]string{e}, p...))
		}
	}
	return
}
