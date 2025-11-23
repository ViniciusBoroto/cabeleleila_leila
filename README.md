# Cabeleleila Leila

Sistema simples de gerenciamento de salão de beleza.

## Tecnologias

### Backend
- Go
- Gin
- GORM
- SQLite
- JWT

### Frontend
- React
- TypeScript
- Tailwind CSS

---

# 📷 Showcase
Showcase do projeto está na pasta `showcase/`.
Inclui:
- Screenshots do programa em execução
- Vídeos demonstrando o funcionamento do programa

---

# 🚀 Como executar
## Requisitos
- [Node](https://nodejs.org/en/download)
- [Go](https://go.dev/doc/install)

## 1. Clonar o repositório
```bash
git clone https://github.com/ViniciusBoroto/cabeleleila_leila.git
cd cabeleleila_leila
```

---

# 🖥️ Frontend
```bash
cd web
npm install
npm run dev
```
O frotnend iniciará em:  
👉 http://localhost:5173/

---

# 🔧 Backend
```bash
cd server
go run main.go
```
O backend iniciará em:  
👉 http://localhost:8080

---

# 🎯 Observações

- O banco SQLite será criado automaticamente no diretório `server/`.
- O frontend e backend funcionam de forma independente. Cada um precisa ser executado em um terminal diferente.
- Ficou pendente a implementação de retry e refresh token pela parte do frontend. No backend já está implementado.
