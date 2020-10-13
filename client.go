package main

import (
	"fmt"
	"bufio"
	//"io"
	"net"
	"strings"
	"time"
	"os"
)

func main() {
	fmt.Println("On se connecte ...")

// --- le client se connecte au port TCP - il tente pendant 1 seconde
	conn, err := net.DialTimeout("tcp", "134.214.202.45:8000", time.Duration(1*time.Second))
    if err != nil {
        fmt.Print("il y a une erreur. Connexion refusée.")
        os.Exit(3)
	}
	
	//ouvre une go routine pour lire les messages reçu par le serveur
	go ecriremessages(conn)
    
    //dans le main: on lis les messages du Terminal:
    
    //1- On digite notre username et on l'envoie
	fmt.Print("Donne moi ton ptit nom: ")
	nom, err:= bufio.NewReader(os.Stdin).ReadString('\n')//bloquante
	check(err)
	conn.Write([]byte("TCCHAT_REGISTER\t" + nom ))
    
    //2- on digite noter message pour les autres:
	condition := true
    for condition{//on devrais arreter la boulce si on se deconnecte
		fmt.Print("Enter text: ")
		texte, err := bufio.NewReader(os.Stdin).ReadString('\n')
		check(err)
		//fmt.Println(texte)
		if (strings.Split(texte, "\n")[0] == "exit") {
			conn.Write([]byte("TCCHAT_DISCONNECT"+"\n"))
            		condition=false

		} else {
			conn.Write([]byte("TCCHAT_MESSAGE\t" + texte ))
		}
	}
}

func identify(ligne string) string {	
	
	array := strings.Split(ligne, "\t")
	message :=" "

	if array[0] == "TCCHAT_WELCOME" {
		mess := array[1]
        	message = "Bienvenue sur "+ mess +" !"+ "\n"
		//fmt.Print("Le serveur "+mess+" vous souhaite la bienvenue"+ "\n")
	}

    //quand le serveur nous previent d'un nouveau arrivé
	if array[0] == "TCCHAT_USERIN" {
		nickname := array[1]
		message = "..... "+nickname+" s'est connecté(e)"+ " ....." +"\n"
		//fmt.Print("=> "+nickname+" s'est deconnecte"+ "\n")
	}
    
    //quand serveur nous previent que qlq est parti
	if array[0] == "TCCHAT_USEROUT" {
		nickname := array[1]
		message = "..... "+nickname+" s'est déconnecté. Au revoir "+ nickname + "!"+ "\n"
		//fmt.Print("=> "+nickname+" s'est deconnecte"+ "\n")
	}
    
    //quand serveur nous broadcast un message de quelqu'un
	if array[0] == "TCCHAT_BCAST" {
        nickname := array[1]
        mess := array[2]
        message = nickname+": "+mess+"\n"
		//fmt.Print(":: "+nickname+" ::"+message)
	} 
	return message
}

func ecriremessages(conn net.Conn) {
	// creation d'un fichier qui stocke les message d'une connection
	fichier, err := os.Create("messenger.txt")
	check(err)

	for {
		// ecoute d'une reponse - lit quand une ligne arrive
	   	ligne, err := bufio.NewReader(conn).ReadString('\n')//bloquant
        if err != nil {
			fmt.Print("il y a une erreur ici.")
			os.Exit(3)
        }
            
		// enleve les balises, retourne le contenu du message
		message := identify(strings.Split(ligne, "\n")[0]) 
		//debugage: fmt.Print("Message from server: "+message)
        
		// on ecrit dans le fichier
		fichier.WriteString(message)
        //fmt.Printf("on a ecrit dans le fichier : ", message)
	}
}

func check(e error) {
    if e != nil {
        fmt.Print("il y a une erreur.")
        os.Exit(3)
    }
}
