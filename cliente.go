package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"os"
)

type User struct {
	Nome, Sobrenome string
}

type Banco string

func main() {
	if os.Args[1] == "" {
		fmt.Println("Insira o nome da pessoa a ser cadastrada:")
		os.Exit(1)
	}

	if os.Args[2] == "" {
		fmt.Println("Insira o sobrenome da pessoa a ser cadastrada:")
		os.Exit(1)
	}

	args := &User{os.Args[1], os.Args[2]}
	var resposta string

	cliente, erro := jsonrpc.Dial("tcp", "localhost:1234")
	if erro != nil {
		log.Fatal("Erro dialing: ", erro)
		os.Exit(1)
	}
	erro = cliente.Call("Cadastro.CriaBanco", new(Banco), &resposta)
	if erro != nil {
		log.Fatal("Erro invocacao CriaBanco: ", erro)
		os.Exit(1)
	}
	fmt.Printf("Criação do banco: %s\n", resposta)

	erro = cliente.Call("Cadastro.AdicionaUsuario", args, &resposta)
	if erro != nil {
		log.Fatal("Erro invocacao AdicionaUsuario: ", erro)
		os.Exit(1)
	}
	fmt.Printf("Resposta do servidor: %s\n", resposta)
}
