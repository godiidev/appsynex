// src/server.ts
import app from './app';
import { config } from './config/config';
import { testConnection } from './config/database';

const startServer = async () => {
  try {
    // Test database connection
    await testConnection();

    // Start server
    app.listen(config.port, () => {
      console.log(`Server is running on port ${config.port}`);
      console.log(`Environment: ${config.nodeEnv}`);
    });
  } catch (error) {
    console.error('Unable to start server:', error);
    process.exit(1);
  }
};

startServer();