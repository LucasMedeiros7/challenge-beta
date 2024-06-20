package models

type Order struct {
	PedidoId  int          `json:"pedidoId"`
	ClienteId int          `json:"clienteId"`
	Status    string       `json:"status"`
	Itens     []ItemPedido `json:"itens"`
}

type ItemPedido struct {
	ItemId     int `json:"itemId"`
	ProdutoId  int `json:"produtoId"`
	Quantidade int `json:"quantidade"`
}
