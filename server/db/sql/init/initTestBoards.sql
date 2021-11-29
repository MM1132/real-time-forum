INSERT INTO boards (boardID, parentID, name, description, "order", isGroup)
VALUES ( 3, 1, 'Animals', 'All the different animals', 1, 0),
       ( 4, 3, 'Cats', 'Please don''t post anything other than just cats here :)', 2, 0),
       ( 5, 3, 'Dogs', 'This is a strictly NO CATS zone!', 1, 0),
       ( 6, 9, 'Parrots', 'Caw caw!', 0, 0),
       ( 7, 9, 'Ducks', 'Quack quack!', 0, 0),
       ( 8, 10, 'Penguins', 'Yes, penguins are birds', 0, 0),
       ( 9, 3, 'Birds', '', 1, 1),
       (10, 9, 'Flightless', '', 0, 1);
