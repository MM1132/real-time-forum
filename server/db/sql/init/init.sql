INSERT INTO boards (boardID, name)
VALUES (0, 'Forum');

INSERT INTO boards (parentID, name, description, "order")
VALUES (0, 'Main Forum', 'The primary place', 1);

INSERT INTO boards (parentID, name, description, "order")
VALUES (0, 'Off-Topic', 'The secondary place', 2);
