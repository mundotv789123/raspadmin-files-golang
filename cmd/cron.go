package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var cronCommand = cli.Command{
	Name:   "cron",
	Usage:  "Start the cron job",
	Action: runCron,
}

func runCron(_ context.Context, cmd *cli.Command) error {
	//TODO: implementar geração de thumbnails e limpeza de arquivos temporários

	/*
		lista todos os arquivos em cada pasta
		pegar todos os arquivos da pasta salvas no banco de dados
		  compara data criacao e alteracao e tamanho do arquivo e arquivos que n foram encontrados
		  adiciona ou atualiza os arquivos
		  chama o serviço de geração de thumbnail para os arquivos que n tem thumbnail ou que foram alterados

		quando for encontrar arquivo no banco que n existe na pasta, apaga thumbnail e remove do banco de dados
	*/
	return nil
}
