// src/config/config.ts
import dotenv from 'dotenv';

dotenv.config();

export const config = {
    port: process.env.PORT || 3000,
    nodeEnv: process.env.NODE_ENV || 'development',  // Sửa nodeenv thành nodeEnv
    jwt: {
        secret: process.env.JWT_SECRET || 'your_secret_key',
        expiresIn: process.env.JWT_EXPIRES_IN || '24h',  // Sửa expiration thành expiresIn
    },
    database: {
        host: process.env.DB_HOST || 'localhost',
        port: Number(process.env.DB_PORT) || 3306,
        username: process.env.DB_USER || 'root',
        password: process.env.DB_PASS || '',
        database: process.env.DB_NAME || 'appsynex',
        dialect: 'mysql' as const,
    },
};