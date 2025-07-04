export interface UserAttributes {
    id?: number;
    username: string;
    password_hash: string;
    email: string;
    phone?: string;
    last_login?: Date;
    account_status: 'active' | 'inactive' | 'suspended';
  }
  
  export interface CreateUserDto {
    username: string;
    password: string;
    email: string;
    phone?: string;
    roleIds: number[];
  }