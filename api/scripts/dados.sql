insert into usuarios (nome, nick, email, senha)
values 
('usuario 1', 'usu01', 'usuario01@gmail.com', '$2a$10$cBOrv3ZSLwTLYGUHCHvBZ.zXeYJLIuAACrM/sayxExwQ45So0yKH2'),
('usuario 2', 'usu02', 'usuario02@gmail.com', '$2a$10$cBOrv3ZSLwTLYGUHCHvBZ.zXeYJLIuAACrM/sayxExwQ45So0yKH2'),
('usuario 3', 'usu03', 'usuario03@gmail.com', '$2a$10$cBOrv3ZSLwTLYGUHCHvBZ.zXeYJLIuAACrM/sayxExwQ45So0yKH2'),
('usuario 4', 'usu04', 'usuario04@gmail.com', '$2a$10$cBOrv3ZSLwTLYGUHCHvBZ.zXeYJLIuAACrM/sayxExwQ45So0yKH2'),
('usuario 5', 'usu05', 'usuario05@gmail.com', '$2a$10$cBOrv3ZSLwTLYGUHCHvBZ.zXeYJLIuAACrM/sayxExwQ45So0yKH2'),
('usuario 6', 'usu06', 'usuario06@gmail.com', '$2a$10$cBOrv3ZSLwTLYGUHCHvBZ.zXeYJLIuAACrM/sayxExwQ45So0yKH2'),
('usuario 7', 'usu07', 'usuario07@gmail.com', '$2a$10$cBOrv3ZSLwTLYGUHCHvBZ.zXeYJLIuAACrM/sayxExwQ45So0yKH2'),
('usuario 8', 'usu08', 'usuario08@gmail.com', '$2a$10$cBOrv3ZSLwTLYGUHCHvBZ.zXeYJLIuAACrM/sayxExwQ45So0yKH2');


insert into seguidores (usuario_id, seguidor_id)
values 
(1,2),
(1,3),
(2,1),
(2,3),
(2,4),
(2,5),
(6,7),
(5,6);

insert into publicacoes (titulo, conteudo, autor_id)
values
("Publicação de usuário 1", "Essa é a publicação do usuário 1! Oba!", 1),
("Publicação de usuário 2", "Essa é a publicação do usuário 2! Oba!", 2),
("Publicação de usuário 3", "Essa é a publicação do usuário 3! Oba!", 3);