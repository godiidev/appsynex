// src/models/Product.ts
@Table({
    tableName: 'products',
    timestamps: true,
    createdAt: 'created_at',
    updatedAt: 'updated_at'
  })
  export class Product extends Model {
    @Column({
      type: DataType.INTEGER,
      primaryKey: true,
      autoIncrement: true
    })
    id!: number;
  
    @ForeignKey(() => ProductName)
    @Column(DataType.INTEGER)
    product_name_id!: number;
  
    @ForeignKey(() => ProductCategory)
    @Column(DataType.INTEGER)
    category_id!: number;
  
    @Column({
      type: DataType.STRING,
      unique: true
    })
    sku!: string;
  
    @Column(DataType.STRING)
    sku_variant?: string;
  
    @Column(DataType.TEXT)
    description?: string;
  
    @Column(DataType.STRING)
    fabric_type?: string;
  
    @Column(DataType.DECIMAL(10, 2))
    weight?: number;
  
    @Column(DataType.DECIMAL(10, 2))
    width?: number;
  
    @Column(DataType.STRING)
    color?: string;
  
    @Column(DataType.STRING)
    quality?: string;
  
    @Column(DataType.STRING)
    fiber_content?: string;
  
    @Column(DataType.JSON)
    additional_info?: object;
  
    @Column(DataType.DECIMAL(10, 2))
    price?: number;
  
    @Column(DataType.DECIMAL(10, 2))
    sales_price?: number;
  
    @Column(DataType.DECIMAL(10, 2))
    stock_quantity?: number;
  
    @BelongsTo(() => ProductName)
    productName!: ProductName;
  
    @BelongsTo(() => ProductCategory)
    category!: ProductCategory;
  }