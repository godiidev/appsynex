//src/models/User.ts
import { Table, Column, Model, DataType, HasMany, BelongsToMany, PrimaryKey, AutoIncrement } from "sequelize-typescript";
import { Role } from "./Role";
import { UserRole } from "./UserRole";
import { col } from "sequelize";

@Table({
    tableName: 'users',
    timestamps: true,
    createdAt: 'created_at',
    updatetAt: 'updated_at'
})

export class User extends Model {
    @Column({
        type: DataType.INTEGER,
        PrimaryKey: true,
        AutoIncrement: true
    })
    id!: number;

    @column({
        type: DataType.STRING,
        allowNull: false,
        unique: true
    })
    username!: string;

    @Column({
        type: DataType.STRING,
        allowNull: false
    })
    password_hash!: string;

    @Column(DataType.STRING)
    email!: string;

    @Column(DataType.STRING)
    phone?: string;

    @Column(DataType.DATE)
    last_login?: Date;

    @Column({
    type: DataType.STRING,
    defaultValue: 'active'
    })
    account_status!: string;

    @BelongsToMany(() => Role, () => UserRole)
    roles?: Role[];
}