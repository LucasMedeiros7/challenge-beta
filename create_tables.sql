-- Criação da tabela Pedido
CREATE TABLE IF NOT EXISTS Pedidos (
    pedidoId SERIAL PRIMARY KEY,
    clienteId INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL
);

-- Criação da tabela ItemPedido
CREATE TABLE IF NOT EXISTS ItensPedido (
    itemId SERIAL PRIMARY KEY,
    pedidoId INTEGER NOT NULL,
    produtoId INTEGER NOT NULL,
    quantidade INTEGER NOT NULL,
    CONSTRAINT fk_pedido FOREIGN KEY (pedidoId) REFERENCES pedidos(pedidoId)
);
