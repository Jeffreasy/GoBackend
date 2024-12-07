
# 🏃‍♂️ DeKoninklijke Loop 2025 - Backend Documentatie

## 📋 Inhoudsopgave
- Systeem Architectuur
- Database Schema
- API Endpoints
- Configuratie
- Beveiliging
- Deployment
- Ontwikkeling
- Troubleshooting

---

## 🏗 Systeem Architectuur
- **Backend**: Go (Golang) 1.22
- **Database**: PostgreSQL
- **Deployment**: Docker + Render
- **Security**: JWT + SSL/TLS
- **Email**: SMTP

---

## 💾 Database Schema

### 👤 Users Tabel
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL
);
```

### 📧 Contacts Tabel
```sql
CREATE TABLE contacts (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    message TEXT NOT NULL
);
```

### 📝 Registrations Tabel
```sql
CREATE TYPE user_role AS ENUM ('deelnemer', 'begeleider', 'vrijwilliger');
CREATE TYPE distance AS ENUM ('2.5km', '6km', '10km', '15km');
CREATE TYPE support_type AS ENUM ('ja', 'nee', 'anders');

CREATE TABLE registrations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'deelnemer',
    distance distance NOT NULL DEFAULT '2.5km',
    needs_support support_type NOT NULL DEFAULT 'nee',
    support_details TEXT,
    terms_accepted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

---

## 🔌 API Endpoints

### 🔐 Authentication
#### **POST** `/register`
**Request:**
```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securePassword123"
}
```

#### **POST** `/login`
**Request:**
```json
{
    "email": "john@example.com",
    "password": "securePassword123"
}
```

### 📨 Contact
#### **POST** `/contact`
**Request:**
```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "message": "Dit is een test bericht"
}
```

### 📋 Registrations
#### **POST** `/registrations`
**Headers:**
```
Authorization: Bearer <jwt_token>
```
**Request:**
```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "role": "deelnemer",
    "distance": "10km",
    "needs_support": "nee",
    "terms_accepted": true
}
```

#### **GET** `/registrations`
**Headers:**
```
Authorization: Bearer <jwt_token>
```

### 🏥 Health Check
#### **GET** `/health`
**Response:**
```json
{"status": "healthy"}
```

---

## ⚙️ Configuratie

### **Server**
```env
SERVER_PORT=8080
```

### **Database**
```env
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-user
DB_PASSWORD=your-password
DB_NAME=your-db-name
```

### **Security**
```env
JWT_SECRET=your-secret-key
```

### **Email**
```env
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your-smtp-user
SMTP_PASSWORD=your-smtp-password
FROM_EMAIL=no-reply@example.com
```

---

## 🔒 Beveiliging

### Allowed Origins
```json
["https://*", "http://*"]
```

### Allowed Methods
```json
["GET", "POST", "PUT", "DELETE", "OPTIONS"]
```

### Allowed Headers
```json
["Accept", "Authorization", "Content-Type"]
```

### Max Age
```text
300 seconds
```

### 🚦 Rate Limiting
- **Max requests**: 100/minuut/IP
- **Burst**: 5
- **Cooldown**: 1 minuut

### 🔑 JWT Authentication
- **Token Expiratie**: 24 uur
- **Algorithm**: HS256
- **Claims**: `user_id`, `exp`

---

## 🚀 Deployment

### 📦 Docker
#### Build
```bash
docker build -t dkl-backend .
docker run -p 8080:8080 dkl-backend
```

#### Compose
```bash
docker-compose up
```

### ☁️ Render
1. Connect GitHub repository
2. Set environment variables
3. Enable auto-deploy

---

## 💻 Ontwikkeling

### 🛠️ Prerequisites
- **Go** 1.22+
- **Docker**
- **PostgreSQL**
- **Git**

### 🏃‍♂️ Running Locally
```bash
# Clone repository
git clone https://github.com/Jeffreasy/GoBackend

# Setup environment
cp .env.example .env

# Start database
docker-compose up db

# Run application
go run cmd/api/main.go
```

---

## ❗ Troubleshooting

### 🔍 Common Issues

#### 1. Database Connection
- Check connection string
- Verify SSL mode
- Test network connectivity

#### Migrations
- Check migration logs
- Verify file order
- Manual database inspection

#### Authentication
- Verify JWT token
- Check expiration
- Validate credentials

### 📊 Logging
```go
log.Printf("Error: %v", err)
```

### 🩺 Diagnostics
```sql
-- Check tables
SELECT * FROM information_schema.tables 
WHERE table_schema = 'public';

-- Check ENUM types
SELECT * FROM pg_type 
WHERE typname IN ('user_role', 'distance', 'support_type');
```
