describe('ThreadTalk - Login Page Smoke Test', () => {

  // CORE LOGIN FUNCTIONALITY TESTS
  it('loads login page, fills the form, toggles password visibility, and enables submit', () => {
    cy.visit('/login');

    // Page renders
    cy.get('.login-page').should('be.visible');
    cy.contains('Welcome Back').should('be.visible');
    cy.contains('Log in to Vitilo ThreadTalk').should('be.visible');

    // Form inputs exist
    cy.get('input[formcontrolname="email"]').should('be.visible');
    cy.get('input[formcontrolname="password"]').should('be.visible');

    // Initially disabled because form invalid
    cy.contains('button[type="submit"]', 'Log In').should('be.disabled');

    // Fill email + password (>= 8 chars to satisfy minlength)
    cy.get('input[formcontrolname="email"]')
      .type('test@example.com')
      .should('have.value', 'test@example.com');

    cy.get('input[formcontrolname="password"]')
      .type('password123')
      .should('have.value', 'password123');

    // Toggle password visibility
    cy.get('button[mat-icon-button][matSuffix]').click();

    // After valid inputs, submit should be enabled (unless isLoading is true)
    cy.contains('button[type="submit"]', 'Log In').should('not.be.disabled');

    // Optional: click submit (we don't assert login success)
    cy.contains('button[type="submit"]', 'Log In').click();

    // Still should be on login OR show some error; just ensure app didn't crash
    cy.get('.login-page').should('be.visible');
  });

  it('navigates to register using "Create Account" button', () => {
    cy.visit('/login');
    cy.contains('button', 'Create Account').click();
    cy.url().should('include', '/register');
  });
  it('submit button has correct disabled state with invalid email', () => {
    cy.visit('/login');
    
    // Fill with invalid email
    cy.get('input[formcontrolname="email"]')
      .type('invalid-email')
      .should('have.value', 'invalid-email');

    cy.get('input[formcontrolname="password"]')
      .type('password123')
      .should('have.value', 'password123');

    // Submit button should remain disabled with invalid email
    cy.contains('button[type="submit"]', 'Log In').should('be.disabled');
  });

  it('submit button has correct disabled state with short password', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]')
      .type('test@example.com')
      .should('have.value', 'test@example.com');

    // Password less than 8 characters
    cy.get('input[formcontrolname="password"]')
      .type('short')
      .should('have.value', 'short');

    // Submit button should be disabled
    cy.contains('button[type="submit"]', 'Log In').should('be.disabled');
  });

  it('"Back to Home" link is clickable and has correct href', () => {
    cy.visit('/login');
    
    const backToHomeLink = cy.contains('a', 'Back to Home');
    
    // Link should be visible
    backToHomeLink.should('be.visible');
    
    // Link should have href attribute pointing to home
    backToHomeLink.should('have.attr', 'routerLink', '/');
  });


  it('navigates back to home using "Back to Home" link', () => {
    cy.visit('/login');
    cy.contains('a', 'Back to Home').click();
    cy.url().should('eq', `${Cypress.config().baseUrl}/`);
  });
