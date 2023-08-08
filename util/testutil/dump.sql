INSERT INTO 
    tasks (title, description, user_id, deadline, 
        created_at, completed_at) 
    VALUES ('full task', 'description', 1, 1691998920, 1691491320, 1691492320);

INSERT INTO 
    tasks (title, user_id, deadline, 
        created_at, completed_at) 
    VALUES ('task without description', 1, 1691998920, 1691491320, 1691492320);

INSERT INTO 
    tasks (title, description, user_id, 
        created_at, completed_at) 
    VALUES ('task without deadline', 'description', 1, 1691491320, 1691492320);

INSERT INTO 
    tasks (title, description, user_id, deadline, 
        created_at, completed_at) 
    VALUES ('full task 2', 'description', 1, 1691998920, 1691491320, 1691492320);

INSERT INTO 
    tasks (title, description, user_id, deadline, 
        created_at) 
    VALUES ('uncompleted task', 'description', 1, 1691998920, 1691491319);

INSERT INTO 
    tasks (title, description, user_id, deadline, 
        created_at) 
    VALUES ('bugged task', NULL, -1, 1691998920, 1691491319);

INSERT INTO 
    tasks (title, description, user_id, deadline, 
        created_at) 
    VALUES ('uncompleted task', 'description', -1, 1691998920, 1691491319);
INSERT INTO 
    tasks (title, description, user_id, deadline, 
        created_at) 
    VALUES ('uncompleted task 3', 'description', 1, 1691998920, 1691491319);