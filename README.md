# Structure of application

## Proposta
Cada histório, deve conter o título, conteúdo, fluxo, campos únicos

### Authentication
**Título**: LOGIN SOCIAL (Primeiro login)
**Conteúdo**: O usuário irá fazer o login social, e caso seja o primeiro login, já criar o tenant e redirecionar para a tela de selecionar o plano
**Fluxo**:  login > autenticação social > criar o usuário> criar o tenant > selecionar o plano > pagina inicial

**Título**:  LOGIN SOCIAL
**Conteúdo**: Quando o usuário for logar, validar se a conta existe, se existe validar se o mesmo tem o tenant e se sim, redicionar para o dashboard inicial
**Fluxo**: > login > login social > validar usuário existe > validar tenant > pagina inicial

**Título**: Check Login
**Conteúdo**: Validar se o usuário está logado e o token é valido. Necessário abrir o token e ver o tempo de expiração
**Fluxo**:  login > validar token se expirado

**Título**: Logout
**Conteúdo**: Validar se o usuário está logado e o token é valido. Necessário abrir o token e ver o tempo de expiração.  Revogar token
**Fluxo**: logout > validar o token se é valido > revoke token

### USER
**Título**: Registro de nov usuário
**Conteúdo**: O usuário irá se registrar no sistema ou logar via Social media.  Ao fazer o registro do usuário, cria automaticamente o tenant.  Após criar o tenant é necessário selecionar o tipo de plano irá seguir.
**Fluxo**: Create user > create tenant >  choose the plan
Campos únicos: user email, tenant name, plan name

**Título**: Login do usuário
**Conteúdo**: O usuário, quando fazer o login, é necessário validar para ver se o usuário existe e se o mesmo pertence a um tenant.  
	**Case**:
		- Se o usuário não existir, retornar erro de usuário não existe;
		- Se o usuário existir e o tenant exisitr, mas a senha não coencidir, retornar erro de invalid login
		- Se o usuário existir e não pertencer a nenhum tenant, retornar erro de tenant não existe;
		- Se tudo for OK, autenticar o usuário e gerar o token
**Fluxo**: Login user > validate password > Validate tenant > gerar token

**Título**: Editar o usuário
**Conteúdo**: Apenas o usuário (owner) pode editar seu próprio profile e o mesmo precisa estar autenticado
	**Case**:
		- Validar se o usuário está autenticado, se não autenticado, retorna erro;
		- Pegar o token JWT do usuário, através de header, se não existir token, retorna erro;
		- Dentro do token, pegar o email do usuário, se o token expirou ou não tiver email válido, retorna erro;
		- Validar se o email do token é o mesmo e-mail do profile a ser atualizado, caso contrário retornar unauthorized;
		- Se for atualizar email e telfone, encaminhar um código de validação, para permitir atualizar;
		- Email do token sendo igual ao do profile a ser atualizado, pode atualizar o profile;
**Fluxo**:  Validar se o usuário está logado > validar o token > match entre token e profile > atualizar o profile

**Título**: Deletar o usuário
**Conteúdo**: Apenas o usuário (owner) pode deletar seu próprio profile e o mesmo precisa estar autenticado
	**Case**:
		- Validar se o usuário está autenticado, se não autenticado, retorna erro;
		- Pegar o token JWT do usuário, através de header, se não existir token, retorna erro;
		- Dentro do token, pegar o email do usuário, se o token expirou ou não tiver email válido, retorna erro;
		- Validar se o email do token é o mesmo e-mail do profile a ser deletado, caso contrário retornar unauthorized;
		- Se for atualizar email e telfone, encaminhar um código de validação, para permitir o delete;
		- Email do token sendo igual ao do profile a ser deletado, pode deletar o profile;
**Fluxo**:  Validar se o usuário está logado > validar o token > match entre token e profile > deletar o profile

**Título**: Listar todos os usuários
**Conteúdo**: Apenas os usuários com roles de system_admin=true podem listar todos os usuários, mas sem permissão de editar o mesmo
	**Case**:
		- Validar se o usuário está autenticado, se não autenticado, retorna erro;
		- Pegar o token JWT do usuário, através de header, se não existir token, retorna erro;
		- Dentro do token, pegar o email do usuário, se o token expirou ou não tiver email válido, retorna erro;
		- Validar se o email do token contém a role system_admin=true;
		- Se nào tiver, retorna unauthorized;
		- Se tiver, lista os usuários;
**Fluxo**:  Validar se o usuário está logado > validar o token > match entre token e a role > listar os usuários

## Autenticação
Tenho o Golang, Gin, Goth para autenticação login social, e quero fazer o seguinte processo:
1. autenticar o usuário através do Goth (já está funcionando)
2. Ao receber os dados do usuário, gerar um novo JWT e esse JWT repassar como header Authorization
  2.1 Também criar o method de refresh do token e logout
3. Como passar as roles do usuário para dentro do token
4. Depois de autenticado, em qual momento validar se o token JWT gerado é valido?
5. Como validar a role do usuário no path especifico?

Explique passo a passo e depois mostre um exemplo completo

Golang, estou desenvolvendo uma gerenciado financeiro em Golang, com as seguintes características:

1. Usuário;
	- possui 1 tenant
2. Planos de contrato;
	- Serão 3 planos free, medium e enterprise.
	  Onde o plano free o usuário não pode compartilhar a wallet e os demais planos, pode
3. Tenant:
	- possui 1 owner;
	- possui 1 plano vinculado;
	- possui ANY wallets
4. Wallet:
  - wallet pertence a um único tenant
	- pode ser compartilhado com ANY tenants
	- o owner do wallet é o mesmo do tenant
	- possui o saldo
5. Category;
  - são categorias para vincular as transações (turismo, educação, finanças,..)
6. Transactions:
  - possui uma categoria
	- pertence a uma única carteira
	- tem tipo (dispesa ou receita)
	- possui o valor
7. O usuário, pode ter acesso a outras wallets, que são compartilhado com ele, através dos tenants

Como criar essa estrutura em golang?

Nas Transaction da wallet, quero que cada novo registro, eu registre no firebase e faça uma trigger para o RabbitMQ, para que processe o valor final do wallet.

Como fazer isso, sendo que em certo momentos, terei muito volume de registros numa mesma wallet?


Ordem de desenolvimento
Schemas
1. Users
2. Plan
3. Tenant
4. WalletCategoryTransaction
5. WalletPermission
6. WalletTransaction
7. Module
8. Role
9. Role Permission
10. Module Settings


type User struct {
	ID              uuid.UUID `json:"id" firestore:"id"`
	Name            string    `json:"name" binding:"required" firestore:"name"`
	Email           string    `json:"email" binding:"required,email" firestore:"email"` // Adicionei validação de email
	AvatarURL       string    `json:"avatarUrl" firestore:"avatarUrl"`
	Provider        string    `json:"provider" firestore:"provider"`
	FirstName       string    `json:"firstName" firestore:"firstName"`
	LastName        string    `json:"lastName" firestore:"lastName"`
	NickName        string    `json:"nickName" firestore:"nickName"`
	Description     string    `json:"description" firestore:"description"`
	Location        string    `json:"location" firestore:"location"`
	CreatedAt       time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt" firestore:"updatedAt"`
}

type UserDetail struct {
	ID            uuid.UUID `json:"id" firestore:"id"`
	UserID        uuid.UUID `json:"userId" firestore:"userId"` // Referência ao usuário
	Address       string    `json:"address" firestore:"address"`
	BirthDate     time.Time `json:"birthDate" firestore:"birthDate"`
	PhoneNumber   string    `json:"phoneNumber" firestore:"phoneNumber"`
	// ... outros detalhes ...
	CreatedAt     time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" firestore:"updatedAt"`
}

// Tenant representa um tenant (cliente) do sistema.
type Tenant struct {
	ID          uuid.UUID `json:"id" firestore:"id"`
	Name        string    `json:"name" binding:"required,min=3,max=200" firestore:"name"`
	Alias       string    `json:"alias" firestore:"alias"`
	OwnerID     uuid.UUID `json:"ownerId" firestore:"ownerId"`
	PlanID      uuid.UUID `json:"planId" firestore:"planId"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updatedAt"`
}

// Wallet representa uma carteira financeira.
type Wallet struct {
	ID              uuid.UUID `json:"id" firestore:"id"`
	Name            string    `json:"name" binding:"required,min=3,max=200" firestore:"name"`
	Description     string    `json:"description" firestore:"description"`
	OwnerID         uuid.UUID `json:"ownerId" firestore:"ownerId"`
	TenantID        uuid.UUID `json:"tenantId" firestore:"tenantId"`  // Tenant ao qual a carteira pertence
	Balance         float64   `json:"balance" firestore:"balance"`    // Saldo atual da carteira
	Currency        string    `json:"currency" firestore:"currency"` // Moeda (ex: "USD", "BRL")
	CreatedAt       time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt" firestore:"updatedAt"`
}

// WalletTransaction representa uma transação em uma carteira.
type WalletTransaction struct {
	ID          uuid.UUID `json:"id" firestore:"id"`
	WalletID    uuid.UUID `json:"walletId" firestore:"walletId"`
	Amount      float64   `json:"amount" firestore:"amount"`
	Type        string    `json:"type" firestore:"type"` // "deposit", "withdrawal", "transfer"
	Description string    `json:"description" firestore:"description"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
}

// WalletCategory representa uma categoria para as transações.
type WalletCategory struct {
	ID          uuid.UUID `json:"id" firestore:"id"`
	Name        string    `json:"name" firestore:"name"`
	Description string    `json:"description" firestore:"description"`
	TenantID    uuid.UUID `json:"tenantId" firestore:"tenantId"` // Permite categorias específicas por tenant
}


// Plan representa um plano de assinatura.
type Plan struct {
	ID          uuid.UUID `json:"id" firestore:"id"`
	Name        string    `json:"name" firestore:"name"`
	Description string    `json:"description" firestore:"description"`
	Features    []string  `json:"features" firestore:"features"` // Lista de recursos
	Price       float64   `json:"price" firestore:"price"`       // Preço mensal
}

// Module representa um módulo do sistema.
type Module struct {
	ID          uuid.UUID `json:"id" firestore:"id"`
	Name        string    `json:"name" binding:"required,min=3,max=200" firestore:"name"`
	Description string    `json:"description" firestore:"description"`
}

// Role representa uma role de acesso.
type Role struct {
	ID          uuid.UUID `json:"id" firestore:"id"`
	Name        string    `json:"name" binding:"required,min=3,max=200" firestore:"name"`
	Description string    `json:"description" firestore:"description"`
}

type PermissionLevel string

const (
	PermissionOwner PermissionLevel = "owner"
	PermissionAdmin PermissionLevel = "admin"
	PermissionEdit  PermissionLevel = "edit"
	PermissionView  PermissionLevel = "view"
)


// RolePermission representa uma permissão de acesso a um módulo.
type RolePermission struct {
	ID          uuid.UUID       `json:"id" firestore:"id"`
	RoleID      uuid.UUID       `json:"roleId" firestore:"roleId"`
	ModuleID    uuid.UUID       `json:"moduleId" firestore:"moduleId"`
	UserID      uuid.UUID       `json:"userId" firestore:"userId"`
	Permission  PermissionLevel `json:"permission" firestore:"permission"`
}

// WalletPermission representa uma permissão de acesso a uma carteira.
type WalletPermission struct {
	ID          uuid.UUID       `json:"id" firestore:"id"`
	UserID      uuid.UUID       `json:"userId" firestore:"userId"`
	WalletID    uuid.UUID       `json:"walletId" firestore:"walletId"`
	Permission  PermissionLevel `json:"permission" firestore:"permission"`
}

// ModuleSetting representa as configurações de um módulo.
type ModuleSetting struct {
	ID          uuid.UUID `json:"id" firestore:"id"`
	ModuleID    uuid.UUID `json:"moduleId" firestore:"moduleId"`
	TenantID    uuid.UUID `json:"tenantId" firestore:"tenantId"`
	UserID      uuid.UUID `json:"userId" firestore:"userId"`
	Key         string    `json:"key" firestore:"key"`
	Value       string    `json:"value" firestore:"value"`
}