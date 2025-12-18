# Scripts de Utilidad

## Crear Usuario Administrador

### Opción 1: Usando el script Go

```bash
go run scripts/create_admin.go admin admin123456
```

### Opción 2: Usando el endpoint bootstrap (si no hay usuarios en la BD)

```bash
curl -X POST http://localhost:8080/users/bootstrap \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123456"
  }'
```

### Opción 3: Directamente en la base de datos (MySQL)

```sql
-- Conectar a MySQL
mysql -u <usuario> -p <nombre_base_datos>

-- Insertar usuario admin
-- Nota: Necesitas generar el hash bcrypt primero usando una herramienta online
-- o usando el método SetPassword del modelo User

INSERT INTO users (username, password_hash, scopes, is_active, created_at, updated_at)
VALUES (
  'admin',
  '$2a$10$XXXXXXXXXXXX...',  -- Hash bcrypt de la contraseña
  '["agents:read","agents:write","workflows:read","workflows:write","steps:read","steps:write","n:read","n:write","users:admin"]',
  true,
  NOW(),
  NOW()
);
```

**Nota:** Para generar el hash bcrypt, puedes usar:
- Un generador online: https://bcrypt-generator.com/
- O usar el script Go que ya tiene la lógica integrada

