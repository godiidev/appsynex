// src/models/ProductCategory.ts
@Table({
    tableName: 'product_categories',
    timestamps: true,
    createdAt: 'created_at',
    updatedAt: 'updated_at'
  })
  export class ProductCategory extends Model {
    @Column({
      type: DataType.INTEGER,
      primaryKey: true,
      autoIncrement: true
    })
    id!: number;
  
    @Column({
      type: DataType.STRING,
      allowNull: false
    })
    category_name!: string;
  
    @Column(DataType.INTEGER)
    parent_category_id?: number;
  
    @Column(DataType.TEXT)
    description?: string;
  
    @HasMany(() => Product)
    products?: Product[];
  
    @HasMany(() => SampleProduct)
    samples?: SampleProduct[];
  
    @BelongsTo(() => ProductCategory, 'parent_category_id')
    parentCategory?: ProductCategory;
  
    @HasMany(() => ProductCategory, 'parent_category_id')
    childCategories?: ProductCategory[];
  }