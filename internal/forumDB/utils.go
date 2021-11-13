package forumDB

import (
	"database/sql"
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
	for _, stmt := range stmtSlice {
		nameSlc := r.FindStringSubmatch(stmt)
		if nameSlc == nil {
			log.Panicf("Could not find a function name in %v:\n%v\n", sqlPath, stmt)
		}

		name := nameSlc[1]
		statements[name], err = db.Prepare(stmt)
		utils.FatalErr(err)
	}

	return statements
}
