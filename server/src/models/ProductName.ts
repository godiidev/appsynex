// src/models/ProductName.ts
import { Table, Column, Model, DataType, HasMany } from 'sequelize-typescript';
import { Product } from './Product';
import { SampleProduct } from './SampleProduct';

@Table({
  tableName: 'product_names',
  timestamps: true,
  createdAt: 'created_at',
  updatedAt: 'updated_at'
})
export class ProductName extends Model {
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
  product_name_vi!: string;

  @Column({
    type: DataType.STRING,
    allowNull: false
  })
  product_name_en!: string;

  @Column({
    type: DataType.STRING,
    unique: true
  })
  sku_parent!: string;

  @HasMany(() => Product)
  products?: Product[];

  @HasMany(() => SampleProduct)
  samples?: SampleProduct[];
}