# Autenticação OAuth 2.0 com Google em Go

Este projeto demonstra a implementação de um fluxo seguro de delegação de acesso, utilizando os padrões OAuth 2.0 e OpenID Connect na linguagem Go com o pacote Goth.

## 🚀 Como funciona o fluxo

Abaixo está o registro passo a passo da configuração e execução do projeto.

### 1. A Interface e a Proteção Inicial
A aplicação possui uma página inicial com o acionador do fluxo de login. Como medida de segurança, qualquer tentativa de login sem as chaves corretas configuradas no ambiente é imediatamente bloqueada pelo Google com o erro `401: invalid_client`.

![Tela inicial da aplicação com botão de Entrar com Google](assets/img/tentativa_login_com_google.png)

![Tela do Google informando erro 401 de client inválido](assets/img/falha_login_google.png)

---

### 2. Configuração do Provedor de Identidade (IdP)
Para habilitar a autenticação, a aplicação foi registrada no Google Cloud Console:
* Foi gerada uma nova credencial selecionando a opção **ID do cliente OAuth**.
* O projeto foi categorizado como um **Aplicativo da Web**.
* A segurança do redirecionamento foi garantida registrando a rota exata de callback da aplicação (`http://localhost:3000/auth/google/callback`) nas URIs autorizadas.

![Menu do Google Cloud para criar ID do cliente OAuth](assets/img/criacao_credencial.png)

![Seleção de Aplicativo da Web no painel](assets/img/selecao_tipo_chave.png)

![Preenchimento das URIs de redirecionamento autorizadas](assets/img/preencher_campos_obrigatorios_criar.png)

---

### 3. Integração de Credenciais
O processo de registro gera um **ID do cliente** e uma **Chave secreta do cliente**. Esses dados são mantidos em segurança no arquivo `.env` da aplicação e injetados no momento da execução para autenticar a nossa API com os servidores do Google.

![Tela exibindo o Client ID e Client Secret gerados](assets/img/copiar_chaves_ou_baixar_json.png)

---

### 4. A Experiência do Usuário (Consentimento)
Com o ambiente configurado, o acionamento do login redireciona o usuário para o ambiente seguro do Google:
* O usuário escolhe qual conta deseja utilizar.
* A tela de consentimento do OAuth é exibida, deixando transparente para o usuário quais escopos de dados (como informações pessoais e e-mail) o aplicativo `app_teste` está solicitando acesso.

![Tela do Google pedindo para escolher qual conta usar](assets/img/escolha_da_conta_acesso.png)

![Tela de consentimento informando quais dados serão acessados](assets/img/permitir_acesso_chave-app_test.png)

---

### 5. Conclusão e Captura de Dados
Ao permitir o acesso, o Google retorna o fluxo para a nossa aplicação junto com o código de autorização. O backend Go troca esse código pelo perfil do usuário, exibindo os dados de identidade com sucesso na tela e nos logs do servidor.

![Aplicação exibindo sucesso no login com os dados de perfil retornados](assets/img/sucesso-acesso-oauth-google.png)

---
### Projeto criado para fins de estudo e desenvolvomento.