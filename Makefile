run-backend:
	cd backend && go run ./cmd/api/main.go

run-frontend:
	cd frontend && npm run dev

run-migration:
	cd backend && go run ./cmd/migrate/migrate.go

run-seeder:
	cd backend && go run ./cmd/seeder/main.go