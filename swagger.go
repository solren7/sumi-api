package main

// @title Sumi API
// @version 0.1.0
// @description Accounting app backend API for auth, categories, transactions, stats, and API keys.
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use the format: Bearer {access_token}
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description Programmatic API key for server-to-server access.
