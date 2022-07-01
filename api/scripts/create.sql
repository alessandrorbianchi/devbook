CREATE DATABASE IF NOT EXISTS devbook;
USE devbook;

DROP TABLE IF EXISTS publicacoes;
DROP TABLE IF EXISTS seguidores;
DROP TABLE IF EXISTS usuarios;
DROP TABLE IF EXISTS mensagens_enviadas;
DROP TABLE IF EXISTS mensagens_recebidas;

CREATE TABLE usuarios(
    id int auto_increment primary key,
    nome varchar(50) not null,
    nick varchar(50) not null unique,
    email varchar(50) not null unique,
    senha varchar(100) not null,
    criadoem timestamp default current_timestamp()
) ENGINE=INNODB;

CREATE TABLE seguidores(
   usuario_id int not null, 
   FOREIGN KEY (usuario_id) 
   REFERENCES usuarios(id)
   ON DELETE CASCADE,

   seguidor_id int not null, 
   FOREIGN KEY (seguidor_id) 
   REFERENCES usuarios(id)
   ON DELETE CASCADE,
  
   primary key (usuario_id, seguidor_id)
) ENGINE=INNODB;

CREATE TABLE publicacoes(
	id int auto_increment primary key,
	titulo varchar(50) not null,
	conteudo varchar(300) not null,

	autor_id int not null,
    FOREIGN KEY (autor_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

	curtidas int default 0,
	criadoem  timestamp default current_timestamp()
) ENGINE=INNODB;

CREATE TABLE mensagens_enviadas(
    id int auto_increment primary key,
    mensagem varchar(300) not null,
    
    remetente_id int not null,
    FOREIGN KEY (remetente_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

    destinatario_id int not null,
    FOREIGN KEY (destinatario_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

     criadoem timestamp default current_timestamp(),
     enviadoem timestamp
) ENGINE=INNODB;

CREATE TABLE mensagens_recebidas(
    id int auto_increment primary key,
    
    mensagem_enviada_id int not null,
    FOREIGN KEY (mensagem_enviada_id)
    REFERENCES mensagens_enviadas(id)
    ON DELETE CASCADE,
    
    mensagem varchar(300) not null,
    
    remetente_id int not null,
    FOREIGN KEY (remetente_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

    destinatario_id int not null,
    FOREIGN KEY (destinatario_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

    recebidoem timestamp default current_timestamp()
) ENGINE=INNODB;