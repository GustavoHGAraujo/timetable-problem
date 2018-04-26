DROP DATABASE IF EXISTS otimizacao;
CREATE DATABASE otimizacao;
USE otimizacao;

DROP TABLE IF EXISTS `Professor`;
CREATE TABLE `Professor` (
  `id_professor` int(11) NOT NULL AUTO_INCREMENT,
  `nome_professor` varchar(30) NOT NULL,
  PRIMARY KEY (`id_professor`)
);

DROP TABLE IF EXISTS `Horario`;
CREATE TABLE `Horario` (
  `id_professor` int(11) NOT NULL,
  `dia_semana` int(11) NOT NULL,
  `posicao` int(11) NOT NULL,
  PRIMARY KEY (`id_professor`,`dia_semana`,`posicao`),
  FOREIGN KEY (`id_professor`) REFERENCES `Professor` (`id_professor`)
);

DROP TABLE IF EXISTS `Disciplina`;
CREATE TABLE `Disciplina` (
  `id_disciplina` INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `nome_disciplina` VARCHAR(30) NOT NULL,
  `id_professor` INTEGER(11) NOT NULL,
  UNIQUE KEY (`nome_disciplina`, `id_professor`),
  FOREIGN KEY (`id_professor`) REFERENCES `Professor` (`id_professor`)
);
