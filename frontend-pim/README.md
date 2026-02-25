# Gondolia PIM Frontend

Product Information Management System Frontend

## Features

- **Dashboard**: Overview of products, categories, and recent changes
- **Products Management**: List, search, filter, and manage products
- **Product Details**: Tabs for master data, prices, categories, attributes, variants, and bundles
- **Product Creation**: Wizard to create new products
- **Categories**: Tree view with drag-and-drop support
- **Authentication**: JWT-based auth system (demo: admin/admin)

## Tech Stack

- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS
- Zustand (state management)
- Lucide Icons

## Development

```bash
npm install
npm run dev
```

Open [http://localhost:3002](http://localhost:3002)

## Docker

```bash
docker build -t gondolia-pim .
docker run -p 3002:3002 gondolia-pim
```

## Environment Variables

See `.env.local.example` for required environment variables.

## API Integration

The PIM frontend connects to the Catalog Service for read operations:
- Products API: `http://catalog:8081/api/v1/products`
- Categories API: `http://catalog:8081/api/v1/categories`

Write operations (create/update/delete) are stubbed and will be implemented when the Catalog Write API is ready.

## Demo Credentials

- Email: `admin`
- Password: `admin`
