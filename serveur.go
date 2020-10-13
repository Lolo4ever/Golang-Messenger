package main

import (
	"fmt"
	//"io"
	"net"
	"strings"
	"os"
	"bufio"
)
func main() {
	fmt.Println("Launching server...")
    fmt.Println("...---...---...---...")
// --- le serveur ecoute le port
 	// listen on a TCP Port - initialisation
	ln, _ := net.Listen("tcp", ":8000")
	
    //creer un canal c pour permetre le broadcast, ainsi toute les routines partagerons le même canal
	c := make(chan map[net.Conn]string /*taille buffer ?*/)
	//tabConn := make([]net.Conn, 0)
	
	// on se met a ecouter le canal
	go broadcast(c)

	for {
		// wait for connection on port 
		conn, _ := ln.Accept()
		fmt.Println("\n"+"Nouvelle connection ..."+"\n")
        
        //chaque nouveau client, une go routine
		go handleRequest(conn,c)
		
		
	}

}

func handleRequest(conn net.Conn, c chan map[net.Conn]string){
	//fmt.Println("On gere les messages ...")
    /* On a un canal c pour pouvoir transmettre à TOUT LE MONDE.
    Dans cette routine on ira que ECRIRE dans le canal et gerer la lecture
    des messages de la connection.
    
    On va creer une autre go routine "envoyer" qui va LIRE dans le canal et
    gerer l'envoye des messages dans la conncetion
    */


    //creation de la go routine "envoyer" qui a pour parametre le canal c et la connection conn
    
    //1- welcome connect message
	newmessage := "TCCHAT_WELCOME\tLe chat de TC\n"
	conn.Write([]byte(newmessage + "\n"))

    //2 - listen for messages
	nickname := " " 
   	paquet := make(map[net.Conn]string)


   	for /*instruction diferant de TTCHAT_DISCONECT*/ {

		// ecoute d'une reponse - lit quand une ligne arrive - bloquant
	   	ligne, err := bufio.NewReader(conn).ReadString('\n')//bloquant, le /n est compris
		//fmt.Println("Ligne recu : "+ligne)
		if err != nil {
				fmt.Print("il y a une erreur ici.")
				os.Exit(3)
		}
		//split le message (2 parties: instruction et nickname(ou parfois message))
		array := strings.Split(ligne, "\t")

		if array[0] == "TCCHAT_REGISTER" {
			nickname = strings.Split(array[1], "\n")[0]
		 	message := "TCCHAT_USERIN\t" + nickname + "\n"
			// on ajoute le message avec commme cle la conn
			paquet[conn] = message
			// on envoie dans le canal la map actualisee
			c <- paquet
		}
		if array[0] == "TCCHAT_MESSAGE" {
			payload := array[1]//contient deja le /n
			message := "TCCHAT_BCAST\t" + nickname +  "\t" + payload
			// on ajoute le message avec commme cle la conn
			paquet[conn] = message
			// on envoie dans le canal la map actualisee
			c <- paquet
		}
		if strings.Split(array[0],"\n")[0]  == "TCCHAT_DISCONNECT"{
			message := "TCCHAT_USEROUT\t" + nickname + "\n"
			// on ajoute le message avec commme cle la conn
			paquet[conn] = message
			// on envoie dans le canal la map actualisee
			c <- paquet

			conn.Close()
			break
		    
		}
	}

}




/*Il y a une routine "envoyer" par connection, mais le canal c'est le
 * même pour TOUT LE MONDE. Donc si une autre connection a un message à envoyer,
 * cette fonction va aussi gerer leurs messages. Mais cette routine ne va qu envoyé 
 * à la connection qui lui a été atribué par paramètre
*/
func broadcast(canal chan map[net.Conn]string /* tabConn []net.Conn */){
    //un slice qui contient les différentes connexions
    tabConn := make([]net.Conn, 0)
    for {
        m := <- canal //bloquant
        
        // on recupere dans le canal la connexion (clé de la map) et le message a envoyer (valeur dans la map)
        //message := 
         
        //pour chaque message dans le canal
        for conn, message := range m {
            array := strings.Split(message, "\t")
            
            //si nouveau connecté: rajoute dans la memoire
            if array[0] == "TCCHAT_USERIN" {
                tabConn = append(tabConn, conn)
            //si logout: on enlève de la memoire
            } else if strings.Split(array[0],"\n")[0] == "TCCHAT_USEROUT"{
                fmt.Println("Quelqu'un a fait LOGOFF")
		indicesup := -1
                for i , _ := range tabConn{
                    if tabConn[i]==conn {
			// on re met toutes les conn avant celle en i (qu'on veut supprimer) et toutes celles apres i
			indicesup = i
                        // it unpacks a slice and passes them as separate arguments to a variadic function.
                    }
                }

		tabConn = append(tabConn[:indicesup], tabConn[indicesup+1:]...)
		
            }
            fmt.Println(tabConn)
            //envoyer a tout les utilisateurs dans la mémoire
            for _ , toutesconn := range tabConn {
                toutesconn.Write([]byte(message))
                //fmt.Println("Envoie la ligne: "+message)
            } 
            delete(m, conn);
        }
    }
}


//obs: on va aussi envoyer à nous même les messages du coup (bien ou pas bien?) bien


/*func check(e error) {
    if e != nil {
        fmt.Print("il y a une erreur.")
        os.Exit(3)
    }
}*/
