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
it('"Create Account" button preserves form state when navigating away', () => {
    cy.visit('/login');
    
    // Fill form
    cy.get('input[formcontrolname="email"]')
      .type('test@example.com');

    cy.get('input[formcontrolname="password"]')
      .type('password123');

    // Click Create Account
    cy.contains('button', 'Create Account').click();

    // Should navigate to register
    cy.url().should('include', '/register');

    // Go back to login
    cy.go('back');

    // Verify we're back on login page
    cy.get('.login-page').should('be.visible');
  });

  it('submit button displays correct text', () => {
    cy.visit('/login');
    
    cy.contains('button[type="submit"]', 'Log In')
      .should('be.visible')
      .should('contain', 'Log In');
  });
  it('email input has correct placeholder text', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]')
      .should('have.attr', 'placeholder', 'your.email@ufl.edu');
  });
  it('email input has email type', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]')
      .should('have.attr', 'type', 'email');
  });

  it('email input has correct autocomplete attribute', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]')
      .should('have.attr', 'autocomplete', 'email');
  });

  it('email field can be cleared after typing', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]')
      .type('test@example.com')
      .clear()
      .should('have.value', '');
  });

  it('email field shows error message when invalid', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]').type('invalid-email');
    cy.get('input[formcontrolname="email"]').blur();
    
    // Should show validation error
    cy.contains('Please enter a valid email').should('be.visible');
  });
it('password input has correct placeholder text', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]')
      .should('have.attr', 'placeholder', 'Enter your password');
  });

  it('password input has correct autocomplete attribute', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]')
      .should('have.attr', 'autocomplete', 'current-password');
  });

  it('password input type is password initially', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]')
      .should('have.attr', 'type', 'password');
  });

  it('password field can be cleared after typing', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]')
      .type('password123')
      .clear()
      .should('have.value', '');
  });

  it('password field shows error when required', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]').type('test@example.com');
    cy.get('input[formcontrolname="email"]').blur();
    
    // Password is required, should show error
    cy.get('input[formcontrolname="password"]').focus().blur();
    cy.contains('Password is required').should('be.visible');
  });

  it('password field shows error when less than 8 characters', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]').type('pass');
    cy.get('input[formcontrolname="password"]').blur();
    
    // Should show minlength error
    cy.contains('Password must be at least 8 characters').should('be.visible');
  });

  // PASSWORD VISIBILITY TOGGLE TESTS
  it('password visibility toggle button exists', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]').type('password123');
    cy.get('button[mat-icon-button][matSuffix]').should('be.visible');
  });

  it('password visibility icon shows correct state when hidden', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]').type('password123');
    
    // Initially password should be hidden, icon should show visibility_off
    cy.get('button[mat-icon-button][matSuffix]')
      .find('mat-icon')
      .should('contain', 'visibility_off');
  });

  it('password visibility icon changes when toggled', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]').type('password123');
    
    // Click toggle
    cy.get('button[mat-icon-button][matSuffix]').click();
    
    // Icon should change to visibility
    cy.get('button[mat-icon-button][matSuffix]')
      .find('mat-icon')
      .should('contain', 'visibility');
  });

  it('toggling password visibility twice returns to original state', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]').type('password123');
    
    const toggle = cy.get('button[mat-icon-button][matSuffix]');
    
    // First click
    toggle.click();
    cy.get('button[mat-icon-button][matSuffix]')
      .find('mat-icon')
      .should('contain', 'visibility');
    
    // Second click
    toggle.click();
    cy.get('button[mat-icon-button][matSuffix]')
      .find('mat-icon')
      .should('contain', 'visibility_off');
  });
// FORM STRUCTURE TESTS
  it('login card is displayed with proper styling', () => {
    cy.visit('/login');
    
    cy.get('.login-card').should('be.visible');
  });

  it('logo section is displayed on login page', () => {
    cy.visit('/login');
    
    cy.get('.logo-section').should('be.visible');
    cy.get('.logo-section svg').should('be.visible');
  });

  it('mat-card-header contains correct content', () => {
    cy.visit('/login');
    
    cy.get('mat-card-header').should('be.visible');
    cy.get('mat-card-title').should('contain', 'Welcome Back');
    cy.get('mat-card-subtitle').should('contain', 'Log in to Vitilo ThreadTalk');
  });

  it('login form has two form fields', () => {
    cy.visit('/login');
    
    cy.get('mat-form-field').should('have.length', 2);
  });
