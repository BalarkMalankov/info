package main

import(
	"fmt"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net"
	"bufio"
	"strings"
	"strconv"
)

//para que no se repita, utilizamos un diccionario
//un mapa es una conexion que guarda elementos [clave]valor
var lib map[int]string = make(map[int]string)




func hostname(name string , port int ) string {
	return fmt.Sprintf("%s:%d", port)
}

func readMsg(con net.Conn) string {
	//creamos un reader
	r := bufio.NewReader(con)
	//creamos la conexion con el reader, y leemos el string hasta le salto de linea, luego retorno el msg
	msg, _ := r.ReadString('\n')
	return string.TrimSpace(msg)
}

func send(port int, name ,msg string) string {
	con, _ := net.Dial("tcp",hostname(name,port))
	defer con.Close()
	fmt.Fprintf(con,msg)
}

//recibimos una direccion 
func sendNR(port int, name ,msg string){
	con, _ := net.Dial("tcp",hostname(name,port))
	defer con.Close()
	fmt.Fprintf(con,msg)
}


func cliAdder(port int){
	//recorre toda la libreta por cada una de las direcciones que tengamos en la lbreta
	for dir := range lib{
		//enviamos un mensaje sin Respuesta
		sendNR(dir + 2, name ,fmt.Sprintf("%d",port)) //+2 pq el servidor +1 es el REGISTRADOR y el +2 es el AGREGADOR
	}
}


func handleRegister(con net.Conn)
{
//cerramos la conexion 
//cuando nos estamos registrando necesitamos leer un dato
	defer con.Close()
	//necesito lleer lso datos name y el puerto
	//leo los datos como string , para ello necesitamos una funcion readMsg
	//convierto a un numero
	port, _ := strconv.Atoi(readMsg(con))

	//invocar al cliente Registraddor que recibe el puerto, para que le avise a los demas nodos de la red que este nuevo nodo con uerto port se ha conectado, 
	cliAdder(port)
	//serializamos la libreta
	jlib, _ := json.Marshal(lib)
	//dsps de avisarle a los demas agreamos a la libreta
	fmt.Fprintln(con,string(jlib)) //jlib es la libreta
	
	
	//dsp de aver obtenido el prto lo agrego 
	lib[port] = "localhost"	// ESTO NECESITA UNA SECCION CRITICA, multiples intentos de escribir un lib al mismo tiempo

}
func servRegister(name string , portbase int ){
	//tiene mas sentido cuando incluyamos nuestra funcion main
	//corren en numero distinto por eso se eutiliza enuna funcion
	ln, _ := net.Listen("tcp", hostname(name,portbase +1 )) 
	defer ln.Close()

	for {
		//crea una conexion
		con, _ := ln.Accept()
		//mandamos una funcion handle
		//escucho un nodo que se quiere registrar lo mando a registrarse 
		go handleRegister()
	}

	
}

//Se comunica on un servidor 
func cliRegister(name string, servport, myport int){
	//este serv va a enviarle un mensaje a un servidor
	//save a quien se va a conectar y le envia el nombre, su credencial
	resp := send(servport + 1, name, fmt.Sprintf("%d", myport))
	//con la rspta necesitamos crear un mapa temporal de tipo entero
	temp := make(map[int]string)
	_ = json.Unmarshal([]byte(resp),&temp)
	for port, na := range temp {
		lib[port] = na
	}
	fmt.Println(lib)
}


//añadimos los nuevos clientes, nuevos nodos
func servAdder(name string, potbase int) {
		//tiene mas sentido cuando incluyamos nuestra funcion main
	//corren en numero distinto por eso se eutiliza enuna funcion
	ln, _ := net.Listen("tcp", hostname(name,portbase + 2 )) 
	defer ln.Close()

	for {
		//crea una conexion
		con, _ := ln.Accept()
		//mandamos una funcion handle
		//escucho un nodo que se quiere registrar lo mando a registrarse 
		go handleAdder(con)
	}
}

func handleAdder(con net.Conn)
{
//cerramos la conexion 
//cuando nos estamos registrando necesitamos leer un dato
	defer con.Close()
	//necesito lleer lso datos name y el puerto
	//leo los datos como string , para ello necesitamos una funcion readMsg
	//convierto a un numero
	port, _ := strconv.Atoi(readMsg(con))
	//dsp de aver obtenido el prto lo agrego 
	lib[port] = "localhost"	// ESTO NECESITA UNA SECCION CRITICA, multiples intentos de escribir un lib al mismo tiempo
	fmt.Println(lib)
}

//cuando el ningun nodo en la red conoce se conecta al REGSISTER, el register le 
//responde con la libreta, 
//este avisa a los demas enviandole un msj al servidor Ader
//uno necesita respuesta y el otro no necesita rpta



func main(){

	name := "localhost"
	port := 0
	fmt.Scanf("%d\n", &port)
	//lanzo a los dos servidores de arriba
	go servRegister(name,port)	

	friendPort := 0
	//solicto a este port que me responda de alguan forma si port es diferente de friendport
	fmt.Scanf("%d\n",&friendPort)
	if port != friendPort{
		//agrego a la libreta al friendport, 
		lib[friendPort] = name
		cliRegister(name, friendPort, port)
	}
	//lo lanzo sin go para que bloquee
	servAdder(name, port)
	
}