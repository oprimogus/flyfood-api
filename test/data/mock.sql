-- Customers
INSERT INTO customer (id, name, last_name, cpf, email, phone, created_at) VALUES
(1, 'João', 'Silva', '12345678901', 'joao.silva@email.com', '11999991111', '2024-01-01'),
(2, 'Maria', 'Santos', '23456789012', 'maria.santos@email.com', '11999992222', '2024-01-01'),
(3, 'Pedro', 'Oliveira', '34567890123', 'pedro.oliveira@email.com', '11999993333', '2024-01-01'),
(4, 'Ana', 'Pereira', '45678901234', 'ana.pereira@email.com', '11999994444', '2024-01-01'),
(5, 'Lucas', 'Costa', '56789012345', 'lucas.costa@email.com', '11999995555', '2024-01-01'),
(6, 'Julia', 'Ferreira', '67890123456', 'julia.ferreira@email.com', '11999996666', '2024-01-01'),
(7, 'Miguel', 'Rodrigues', '78901234567', 'miguel.rodrigues@email.com', '11999997777', '2024-01-01'),
(8, 'Beatriz', 'Almeida', '89012345678', 'beatriz.almeida@email.com', '11999998888', '2024-01-01'),
(9, 'Gabriel', 'Lima', '90123456789', 'gabriel.lima@email.com', '11999999999', '2024-01-01'),
(10, 'Sofia', 'Carvalho', '01234567890', 'sofia.carvalho@email.com', '11999990000', '2024-01-01'),
(11, 'Rafael', 'Martins', '11223344556', 'rafael.martins@email.com', '11988881111', '2024-01-01'),
(12, 'Isabella', 'Souza', '22334455667', 'isabella.souza@email.com', '11988882222', '2024-01-01'),
(13, 'Matheus', 'Gomes', '33445566778', 'matheus.gomes@email.com', '11988883333', '2024-01-01'),
(14, 'Laura', 'Ribeiro', '44556677889', 'laura.ribeiro@email.com', '11988884444', '2024-01-01'),
(15, 'Enzo', 'Fernandes', '55667788990', 'enzo.fernandes@email.com', '11988885555', '2024-01-01'),
(16, 'Valentina', 'Barbosa', '66778899001', 'valentina.barbosa@email.com', '11988886666', '2024-01-01'),
(17, 'Thiago', 'Pinto', '77889900112', 'thiago.pinto@email.com', '11988887777', '2024-01-01'),
(18, 'Clara', 'Moreira', '88990011223', 'clara.moreira@email.com', '11988888888', '2024-01-01'),
(19, 'Davi', 'Castro', '99001122334', 'davi.castro@email.com', '11988889999', '2024-01-01'),
(20, 'Alice', 'Cardoso', '00112233445', 'alice.cardoso@email.com', '11988880000', '2024-01-01');

-- Addresses
INSERT INTO address (customer_id, name, address_line_1, address_line_2, neighborhood, city, state, postal_code, latitude, longitude, country, created_at) VALUES
(1, 'Casa', 'Rua Augusta, 1000', 'Apto 101', 'Consolação', 'São Paulo', 'SP', '01304-001', '-23.550520', '-46.652937', 'Brasil', '2024-01-01'),
(2, 'Casa', 'Rua Oscar Freire, 500', 'Casa', 'Jardins', 'São Paulo', 'SP', '01426-001', '-23.562256', '-46.669749', 'Brasil', '2024-01-01'),
(3, 'Casa', 'Avenida Paulista, 1500', 'Apto 1001', 'Bela Vista', 'São Paulo', 'SP', '01311-200', '-23.564616', '-46.652857', 'Brasil', '2024-01-01'),
(4, 'Casa', 'Rua dos Pinheiros, 800', 'Casa 2', 'Pinheiros', 'São Paulo', 'SP', '05422-001', '-23.566347', '-46.678661', 'Brasil', '2024-01-01'),
(5, 'Casa', 'Rua João Cachoeira, 300', 'Apto 502', 'Itaim Bibi', 'São Paulo', 'SP', '04535-001', '-23.585372', '-46.672806', 'Brasil', '2024-01-01'),
(6, 'Casa', 'Alameda Santos, 700', 'Apto 801', 'Jardim Paulista', 'São Paulo', 'SP', '01419-001', '-23.569651', '-46.648727', 'Brasil', '2024-01-01'),
(7, 'Casa', 'Rua Harmonia, 400', 'Casa', 'Vila Madalena', 'São Paulo', 'SP', '05435-001', '-23.555372', '-46.687831', 'Brasil', '2024-01-01'),
(8, 'Casa', 'Rua Amauri, 200', 'Apto 1501', 'Itaim Bibi', 'São Paulo', 'SP', '01448-001', '-23.581480', '-46.674929', 'Brasil', '2024-01-01'),
(9, 'Casa', 'Rua Teodoro Sampaio, 1000', 'Apto 601', 'Pinheiros', 'São Paulo', 'SP', '05406-000', '-23.558736', '-46.676931', 'Brasil', '2024-01-01'),
(10, 'Casa', 'Rua da Consolação, 2000', 'Apto 1201', 'Consolação', 'São Paulo', 'SP', '01301-100', '-23.555372', '-46.662831', 'Brasil', '2024-01-01'),
(11, 'Casa', 'Rua Joaquim Floriano, 500', 'Apto 401', 'Itaim Bibi', 'São Paulo', 'SP', '04534-001', '-23.583372', '-46.675831', 'Brasil', '2024-01-01'),
(12, 'Casa', 'Rua dos Tamboréus, 300', 'Casa', 'Moema', 'São Paulo', 'SP', '04543-001', '-23.601372', '-46.662831', 'Brasil', '2024-01-01'),
(13, 'Casa', 'Rua Tabapuã, 800', 'Apto 901', 'Itaim Bibi', 'São Paulo', 'SP', '04533-001', '-23.584372', '-46.673831', 'Brasil', '2024-01-01'),
(14, 'Casa', 'Rua Haddock Lobo, 600', 'Apto 701', 'Cerqueira César', 'São Paulo', 'SP', '01414-001', '-23.558372', '-46.662831', 'Brasil', '2024-01-01'),
(15, 'Casa', 'Rua Lisboa, 400', 'Casa', 'Jardins', 'São Paulo', 'SP', '05413-001', '-23.567372', '-46.668831', 'Brasil', '2024-01-01'),
(16, 'Casa', 'Rua Bela Cintra, 900', 'Apto 1101', 'Consolação', 'São Paulo', 'SP', '01415-001', '-23.557372', '-46.662831', 'Brasil', '2024-01-01'),
(17, 'Casa', 'Rua Wisard, 200', 'Casa', 'Vila Madalena', 'São Paulo', 'SP', '05434-001', '-23.553372', '-46.687831', 'Brasil', '2024-01-01'),
(18, 'Casa', 'Rua Gabriel Monteiro, 700', 'Apto 801', 'Jardim Paulistano', 'São Paulo', 'SP', '01441-001', '-23.573372', '-46.682831', 'Brasil', '2024-01-01'),
(19, 'Casa', 'Rua Mourato Coelho, 500', 'Casa', 'Vila Madalena', 'São Paulo', 'SP', '05417-001', '-23.554372', '-46.687831', 'Brasil', '2024-01-01'),
(20, 'Casa', 'Rua João Moura, 300', 'Apto 601', 'Pinheiros', 'São Paulo', 'SP', '05412-001', '-23.559372', '-46.677831', 'Brasil', '2024-01-01');

-- Owners (first 10 customers are store owners)
INSERT INTO owner (id, signature_active, created_at) VALUES
(1, true, '2024-01-01'),
(2, true, '2024-01-01'),
(3, true, '2024-01-01'),
(4, true, '2024-01-01'),
(5, true, '2024-01-01'),
(6, true, '2024-01-01'),
(7, true, '2024-01-01'),
(8, true, '2024-01-01'),
(9, true, '2024-01-01'),
(10, true, '2024-01-01');

-- Stores
INSERT INTO store (id, owner_id, cnpj, name, description, active, phone, score, is_open, type, address_line_1, address_line_2, neighborhood, city, state, postal_code, latitude, longitude, country, created_at) VALUES
('01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', 1, '12345678000101', 'Bar do João', 'Bar tradicional com petiscos', true, '11999991111', 4, true, 'PUB', 'Rua Augusta, 1500', '', 'Consolação', 'São Paulo', 'SP', '01304-001', '-23.550520', '-46.652937', 'Brasil', '2024-01-01'),
('01940ec5-8185-7f13-99bf-6b799548e589', 2, '23456789000102', 'Mercado Maria', 'Mercado de bairro completo', true, '11999992222', 5, true, 'MARKET', 'Rua Oscar Freire, 1000', '', 'Jardins', 'São Paulo', 'SP', '01426-001', '-23.562256', '-46.669749', 'Brasil', '2024-01-01'),
('01940ec6-1c8a-7d1f-ab47-3b11a753d28e', 3, '34567890000103', 'Água Mineral SP', 'Distribuidora de água mineral', true, '11999993333', 4, true, 'CONVENIENCE', 'Avenida Paulista, 2000', '', 'Bela Vista', 'São Paulo', 'SP', '01310-200', '-23.564616', '-46.652857', 'Brasil', '2024-01-01'),
('01940ec6-6330-7071-af1d-6a56faba1bc1', 4, '45678901000104', 'Farmácia Ana', 'Farmácia 24 horas', true, '11999994444', 5, true, 'PHARMACY', 'Rua dos Pinheiros, 1200', '', 'Pinheiros', 'São Paulo', 'SP', '05422-001', '-23.566347', '-46.678661', 'Brasil', '2024-01-01'),
('01940ec6-9857-7035-af61-391ce2783608', 5, '56789012000105', 'Bar do Lucas', 'Bar e restaurante', true, '11999995555', 4, true, 'PUB', 'Rua João Cachoeira, 800', '', 'Itaim Bibi', 'São Paulo', 'SP', '04535-001', '-23.585372', '-46.672806', 'Brasil', '2024-01-01'),
('01940ec6-d42e-798d-adfd-067a8ec9dca3', 6, '67890123000106', 'Conveniência Julia', 'Loja de conveniência 24h', true, '11999996666', 4, true, 'CONVENIENCE', 'Alameda Santos, 1200', '', 'Jardim Paulista', 'São Paulo', 'SP', '01419-001', '-23.569651', '-46.648727', 'Brasil', '2024-01-01'),
('01940ec7-08ff-701f-bc6a-b57e34ea8496', 7, '78901234000107', 'Tabacaria Miguel', 'Tabacaria premium', true, '11999997777', 5, true, 'TOBBACO', 'Rua Harmonia, 800', '', 'Vila Madalena', 'São Paulo', 'SP', '05435-001', '-23.555372', '-46.687831', 'Brasil', '2024-01-01'),
('01940ec7-6fc6-7870-9a85-784c8d34b029', 8, '89012345000108', 'Mercado Bia', 'Supermercado completo', true, '11999998888', 4, true, 'MARKET', 'Rua Amauri, 600', '', 'Itaim Bibi', 'São Paulo', 'SP', '01448-001', '-23.581480', '-46.674929', 'Brasil', '2024-01-01'),
('01940ec7-a547-766a-84b2-a8204ad141f2', 9, '90123456000109', 'Bar do Gabriel', 'Bar e cervejaria artesanal', true, '11999999999', 5, true, 'PUB', 'Rua Teodoro Sampaio, 1500', '', 'Pinheiros', 'São Paulo', 'SP', '05406-000', '-23.558736', '-46.676931', 'Brasil', '2024-01-01'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 10, '01234567000110', 'Mercado Sofia', 'Mercado de produtos naturais', true, '11999990000', 4, true, 'MARKET', 'Rua da Consolação, 2500', '', 'Consolação', 'São Paulo', 'SP', '01301-100', '-23.555372', '-46.662831', 'Brasil', '2024-01-01');

-- Products for Bar do João (PUB)
INSERT INTO product (id, store_id, active_for_sale, promo_active, type, tag, name, description, stock_quantity, score, price, created_at) VALUES
('01940ec7-ed01-7cb5-b610-a473d64c9eb9', '01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', true, false, 'FOOD', 'petiscos', 'Pastéis', 'Mix de pastéis variados', 100, 5, 3000, '2024-01-01'),
('01940ec8-1a32-707c-9c08-ad6373e8d822', '01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', true, false, 'FOOD', 'petiscos', 'Batata Frita', 'Porção de batatas fritas', 100, 4, 2500, '2024-01-01'),
('01940ec8-4180-7d50-92d6-7d8e6ac21bdd', '01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', true, true, 'WATER', 'bebidas', 'Água Mineral', 'Garrafa 500ml', 200, 5, 500, '2024-01-01'),
('01940ec8-6570-70cb-8342-335f470262b8', '01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', true, false, 'FOOD', 'promocao', 'Isca de Frango', 'Porção de iscas de frango', 100, 4, 3500, '2024-01-01'),
('01940ec8-8fef-766b-8705-ae0cf1391eac', '01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', true, true, 'FOOD', 'promocao', 'Espetinho', 'Espetinho de carne', 150, 5, 1500, '2024-01-01');

-- Products for Mercado Maria (MARKET)
INSERT INTO product (id, store_id, active_for_sale, promo_active, type, tag, name, description, stock_quantity, score, price, created_at) VALUES
('22222222-2222-2222-2222-222222222200', '01940ec5-8185-7f13-99bf-6b799548e589', true, false, 'FOOD', 'mercearia', 'Arroz', 'Pacote 5kg', 100, 5, 2500, '2024-01-01'),
('22222222-2222-2222-2222-222222222201', '01940ec5-8185-7f13-99bf-6b799548e589', true, false, 'FOOD', 'mercearia', 'Feijão', 'Pacote 1kg', 150, 4, 800, '2024-01-01'),
('22222222-2222-2222-2222-222222222202', '01940ec5-8185-7f13-99bf-6b799548e589', true, true, 'WATER', 'bebidas', 'Água Mineral', 'Galão 20L', 50, 5, 1500, '2024-01-01'),
('22222222-2222-2222-2222-222222222203', '01940ec5-8185-7f13-99bf-6b799548e589', true, true, 'FOOD', 'promocao', 'Macarrão', 'Pacote 500g', 200, 4, 500, '2024-01-01'),
('22222222-2222-2222-2222-222222222204', '01940ec5-8185-7f13-99bf-6b799548e589', true, false, 'FOOD', 'hortifruti', 'Banana', 'Kg', 100, 5, 600, '2024-01-01');

-- Products for Água Mineral SP (CONVENIENCE)
INSERT INTO product (id, store_id, active_for_sale, promo_active, type, tag, name, description, stock_quantity, score, price, created_at) VALUES
('33333333-3333-3333-3333-333333333300', '01940ec6-1c8a-7d1f-ab47-3b11a753d28e', true, false, 'WATER', 'galao', 'Água Mineral', 'Galão 20L', 200, 5, 1200, '2024-01-01'),
('33333333-3333-3333-3333-333333333301', '01940ec6-1c8a-7d1f-ab47-3b11a753d28e', true, true, 'WATER', 'garrafas', 'Água Mineral', 'Garrafa 1.5L', 300, 4, 400, '2024-01-01'),
('33333333-3333-3333-3333-333333333302', '01940ec6-1c8a-7d1f-ab47-3b11a753d28e', true, false, 'WATER', 'garrafas', 'Água Mineral', 'Garrafa 500ml', 500, 5, 250, '2024-01-01'),
('33333333-3333-3333-3333-333333333303', '01940ec6-1c8a-7d1f-ab47-3b11a753d28e', true, true, 'WATER', 'promocao', 'Kit Água', 'Pack 6 garrafas 1.5L', 100, 5, 2000, '2024-01-01'),
('33333333-3333-3333-3333-333333333304', '01940ec6-1c8a-7d1f-ab47-3b11a753d28e', true, false, 'WATER', 'copos', 'Água Mineral', 'Copo 200ml', 1000, 4, 150, '2024-01-01');

-- Continue with more products for other stores...

-- Store Payment Methods
INSERT INTO store_payment_method (id, payment_method) VALUES
('01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', 'CREDIT'),
('01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', 'DEBIT'),
('01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', 'PIX'),
('01940ec5-8185-7f13-99bf-6b799548e589', 'CREDIT'),
('01940ec5-8185-7f13-99bf-6b799548e589', 'DEBIT'),
('01940ec5-8185-7f13-99bf-6b799548e589', 'PIX'),
('01940ec5-8185-7f13-99bf-6b799548e589', 'CASH'),
('01940ec6-1c8a-7d1f-ab47-3b11a753d28e', 'CREDIT'),
('01940ec6-1c8a-7d1f-ab47-3b11a753d28e', 'DEBIT'),
('01940ec6-1c8a-7d1f-ab47-3b11a753d28e', 'PIX');

-- Store Business Hours
INSERT INTO store_business_hour (id, weekday, open_hour, closing_hour) VALUES
('01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', 1, '16:00', '00:00'),
('01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', 2, '16:00', '00:00'),
('01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', 3, '16:00', '00:00'),
('01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', 4, '16:00', '00:00'),
('01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', 5, '16:00', '02:00'),
('01940ec5-29a7-7ab9-beb9-d6ff5c021ab2', 6, '16:00', '02:00'),
('01940ec5-8185-7f13-99bf-6b799548e589', 1, '08:00', '22:00'),
('01940ec5-8185-7f13-99bf-6b799548e589', 2, '08:00', '22:00'),
('01940ec5-8185-7f13-99bf-6b799548e589', 3, '08:00', '22:00'),
('01940ec5-8185-7f13-99bf-6b799548e589', 4, '08:00', '22:00'),
('01940ec5-8185-7f13-99bf-6b799548e589', 5, '08:00', '22:00'),
('01940ec5-8185-7f13-99bf-6b799548e589', 6, '08:00', '22:00'),
('01940ec5-8185-7f13-99bf-6b799548e589', 7, '08:00', '20:00'),
('01940ec6-1c8a-7d1f-ab47-3b11a753d28e', 1, '08:00', '18:00'),
('01940ec6-1c8a-7d1f-ab47-3b11a753d28e', 2, '08:00', '18:00'),
('01940ec6-1c8a-7d1f-ab47-3b11a753d28e', 3, '08:00', '18:00'),
('01940ec6-1c8a-7d1f-ab47-3b11a753d28e', 4, '08:00', '18:00'),
('01940ec6-1c8a-7d1f-ab47-3b11a753d28e', 5, '08:00', '18:00'),
('01940ec6-1c8a-7d1f-ab47-3b11a753d28e', 6, '08:00', '16:00');