package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	//Para acesso ao database
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Nome, Sobrenome string
}

type Cadastro string
type Banco string

func (t *Cadastro) CriaBanco(arg *Banco, res *string) error {
	// Conectando ao mysql
	db, err := sql.Open("mysql", "root:toor@tcp(localhost:3306)/")
	if err != nil {
		log.Fatal(err)
	}

	// Faz um ping no banco para testar a conexão
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Criando Banco 'rpc' caso não exista")
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS rpc")
	if err != nil {
		*res = "Falha"
		log.Fatal(err)
	}
	fmt.Println("Usando schemma rpc")
	_, err = db.Exec("USE rpc")
	if err != nil {
		*res = "Falha"
		log.Fatal(err)
	}

	fmt.Println("Criando tabela pessoa caso não exita")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS rpc.usuario (
              nome VARCHAR(30) NULL,
              sobrenome VARCHAR(30) NULL)`)
	if err != nil {
		*res = "Falha"
		log.Fatal(err)
	}
	db.Close()
	*res = "Sucesso"
	return nil
}

func (t *Cadastro) AdicionaUsuario(user *User, res *string) error {

	// Conectando ao mysql
	db, err := sql.Open("mysql", "root:toor@tcp(localhost:3306)/")
	if err != nil {
		log.Fatal(err)
	}

	// Faz um ping no banco para testar a conexão
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare(`Insert Into usuario
             (nome, sobrenome)
             Values (?, ?)`)
	defer stmt.Close()

	fmt.Println("Tentando inserir no banco")
	_, err = stmt.Exec((*user).Nome, (*user).Sobrenome) // Insere no banco
	if err != nil {
		*res = "Falha"
		log.Fatal(err)
	}

	fmt.Println("Usuario adicionado: " + (*user).Nome + " " + (*user).Sobrenome)
	*res = "Sucesso"
	db.Close()
	return nil
}

func main() {
	servicoCadastro := new(Cadastro)
	rpc.Register(servicoCadastro)
	porta := "localhost:1234"

	l, e := net.Listen("tcp", porta)
	fmt.Println("Ouvindo na porta", porta)
	if e != nil {
		log.Fatal("erro no listen: ", e)
	}

	for {
		conn, e := l.Accept()
		if e != nil {
			log.Fatal("erro na recepcao da conexao: ", e)
		}
		jsonrpc.ServeConn(conn)
	}
}
