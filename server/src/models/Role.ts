// src/models/Role.ts
@Table({
    tableName: 'roles',
    timestamps: true,
    createdAt: 'created_at',
    updatedAt: 'updated_at'
  })
  export class Role extends Model {
    @Column({
      type: DataType.INTEGER,
      primaryKey: true,
      autoIncrement: true
    })
    id!: number;
  
    @Column(DataType.STRING)
    role_name!: string;
  
    @Column(DataType.TEXT)
    description?: string;
  
    @BelongsToMany(() => User, () => UserRole)
    users?: User[];
  
    @BelongsToMany(() => Permission, () => RolePermission)
    permissions?: Permission[];
  }