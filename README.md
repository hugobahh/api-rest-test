##SQL
CREATE TABLE `Usr` (
  `id_Usr` int(5) NOT NULL AUTO_INCREMENT,
  `User` varchar(100) DEFAULT NULL,
  `Mail` varchar(50) DEFAULT NULL,
  `Pwd` varchar(50) DEFAULT NULL,
  `St` varchar(10) DEFAULT 'ACTIVE',
  `Date_reg` datetime DEFAULT NULL,
  `Tel` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id_Usr`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;

## INTRODUCCIÓN
Este repositorio contiene una serie de requerimientos de un Caso Práctico, que busca evaluar las capacidades técnicas del candidato con respecto a las principales funciones y responsabilidades que se requieren dentro del área de Desarrollo de Tecnología de _GCP Global_.

#### EJERCICIO:
Realiza el siguiente ejercicio, comparte la liga del repositorio de tu prueba. 

1. Crear un servicio de registro de usuario que reciba como parámetros usuario, correo, teléfono y contraseña.

2. El servicio deberá validar que el correo y telefono no se encuentren registrados, de lo contrario deberá retornar un mensaje “el correo/telefono ya se encuentra registrado”.

3. Deberá validar que la contraseña sea de 6 caracteres mínimo y 12 máximo y contener al menos una mayúscula, una minúscula, un carácter especial (@ $ o &) y un número.

4. Validar que el teléfono sea a 10 dígitos y el correo tenga un formato válido.

5. Crear un servicio login que reciba como parámetros usuario o correo y contraseña.
6. El servicio debe devolver un token jwt.

7. Deberá validar que el usuario o correo y contraseña sean válidos, de lo contrario retorna un mensaje “usuario / contraseña incorrectos”.

8. En ambos servicios se deberá validar que todos los parámetros solicitados vayan en el cuerpo de la petición, de lo contrario retorna un mensaje con el campo faltante.

Ejemplo login:
{
“correo”: “prueba@gmail.com”
}

La respuesta debe ser:
“Falta el campo contraseña”

#### RECOMENDACIONES:
1. Usaría redis para gestionar sesiones con el token de jwt.

2. Migrar los servicos que sean posibles para un mejor performance usando Graphql, grpc y protobuffers.
 