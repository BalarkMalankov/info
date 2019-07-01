package main
 
import(
    "bufio"
    "fmt"
    "net"
    "strconv"
    "strings"  
    "time"
)
 
/*NOS INCIAMOS COMO SERVIDOR*/
func server(ch,end chan int){//de pie
    ln, _ := net.Listen("tcp","10.21.61.175:8000")
    defer ln.Close()
    /*el servidor espera por conexiones*/
    con, _ := ln.Accept()
    defer con.Close()
    r := bufio.NewReader(con) //necesitamos un vector con para nuestaa conexion
    //declara una variable par leer lo que se envia por consola
    /*recibe muchos mensajes*/
   
    for{
        /*leemos el mensaje que se nos esta enviando*/
        //que ha recibido el servidor
        msg, _ := r.ReadString('\n')
        //funcion para quitar los espaciados de los mensajes
        msg = strings.TrimSpace(msg)
        fmt.Println("Recibido: ", msg,)
        /*La funcion atoi convierte de texto a numero, y me bota el error, si no hubo errores, el error va a ser nulo*/
        if n, err := strconv.Atoi(msg); err == nil{
            /*si no hay error vamos a ejecutar*/
            if n==0{
                fmt.Println("Recáspita, perdí! 😒")
               
            }else{              /*le damos un dejaly de un 1seg*/
 
                time.Sleep(time.Second)
            }
            /*enviamos al canal un numero -1 (osea el mensaje a continuacion)*/
            ch<- n-1
            if n<0{
                break
            }
        }
 
    }
    /*cuando termina el servidor envia 0 al canal*/
    end<- 0
 
}
 
/*vamos a tener una coneccion*/
func client(ch chan int){
    /*crear una conexion */
    var con net.Conn
    /*para que no se quede eesperando a que selibere el canal*/
    created := false
    for{
        n := <-ch
        /*si creo que la conexion no ha sido creada*/
        if !created {
            created = true
            //va a hacer la conexon 
            //establecer la conexion desde una pc a otra (credenciales de la otra ip)
            con, _ = net.Dial("tcp","10.21.61.168:8000")
        }
        fmt.Println("Enviando: ", n)
        fmt.Fprintf(con,"%d\n", n)
        /*desps de enviar si el numero que hemos recibido es menor que -1 salimos*/
        if n < -1 {
            fmt.Println("Ufffff no Perdí")
 
            break
        }
    }
    defer con.Close()
 
}
 
 
/*con seccion critica*/
func start(ch chan int){
    var n int
    fmt.Scanf("%d\n",&n)
    ch<-n
}
 
 
func main(){
 
    /*el programa va a esperar hasta que el servidor escriba algo al canal*/
   
    ch:= make(chan int)
    end:= make(chan int)
    go server(ch,end)
    go client(ch)
    go start(ch)
   
    <-end
}