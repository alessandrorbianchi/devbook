package seguranca

import (
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func Hash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

func VerificaSenha(senhaComHash, senhaString string) error {
	return bcrypt.CompareHashAndPassword([]byte(senhaComHash), []byte(senhaString))
}

func GerarCodigoDeSeguranca(remetenteID, destinatarioID uint64) (uint64, error) {
	codigoDeSeguranca := fmt.Sprintf("%d%d", remetenteID, destinatarioID)
	if remetenteID > destinatarioID {
		codigoDeSeguranca = fmt.Sprintf("%d%d", destinatarioID, remetenteID)
	}

	codigo, err := strconv.ParseUint(codigoDeSeguranca, 10, 64)
	if err != nil {
		return 0, nil
	}

	return codigo, nil
}
