CREATE USER IF NOT EXISTS 'servico_usuario'@'%' IDENTIFIED BY 'sua_senha_aqui';
GRANT ALL PRIVILEGES ON mydb.* TO 'servico_usuario'@'%';
FLUSH PRIVILEGES;