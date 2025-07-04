// src/routes/v1/index.ts
import { Router } from 'express';

const router = Router();

// Health check route
router.get('/test', (req, res) => {
  res.status(200).json({ message: 'API v1 is working' });
});

export default router;