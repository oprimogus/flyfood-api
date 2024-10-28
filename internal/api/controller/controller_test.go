package controller_test

// import (
// 	"context"
// 	"os"
// 	"testing"

// 	"github.com/oprimogus/cardapiogo/test/integration"
// )

// func TestMain(m *testing.M) {
// 	ctx := context.Background()
// 	postgres, err := integration.MakePostgres(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer postgres.Kill(ctx)

// 	keycloak, err := integration.MakeKeycloak(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer keycloak.Kill(ctx)

// 	code := m.Run()
// 	os.Exit(code)
// }
