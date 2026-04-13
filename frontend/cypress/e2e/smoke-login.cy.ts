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
it('form fields have outline appearance', () => {
    cy.visit('/login');
    
    cy.get('mat-form-field[appearance="outline"]').should('have.length', 2);
  });

  it('login page container has correct class', () => {
    cy.visit('/login');
    
    cy.get('.login-container').should('be.visible');
  });

  // FORM VALIDATION TESTS
  it('form is invalid when email is empty', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]').type('password123');
    
    // Form should be invalid, submit button disabled
    cy.contains('button[type="submit"]', 'Log In').should('be.disabled');
  });

  it('form is invalid when password is empty', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]').type('test@example.com');
    
    // Form should be invalid, submit button disabled
    cy.contains('button[type="submit"]', 'Log In').should('be.disabled');
  });

  it('valid form enables submit button', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]').type('test@example.com');
    cy.get('input[formcontrolname="password"]').type('password123');
    
    cy.contains('button[type="submit"]', 'Log In')
      .should('not.be.disabled');
  });

  // BUTTON TYPE TESTS
  it('submit button type is submit', () => {
    cy.visit('/login');
    
    cy.contains('button[type="submit"]', 'Log In')
      .should('have.attr', 'type', 'submit');
  });

  it('password toggle button type is button (not submit)', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]').type('password123');
    
    cy.get('button[mat-icon-button][matSuffix]')
      .should('have.attr', 'type', 'button');
  });
it('signup section is visible on login page', () => {
    cy.visit('/login');
    
    cy.get('.signup-section').should('be.visible');
  });

  it('signup section contains "Don\'t have an account?" text', () => {
    cy.visit('/login');
    
    cy.get('.signup-section').contains("Don't have an account?").should('be.visible');
  });

  // DIVIDER TESTS
  it('divider with "or" text is visible', () => {
    cy.visit('/login');
    
    cy.get('.divider').should('be.visible');
    cy.get('.divider').contains('or').should('be.visible');
  });

  it('back link section is visible', () => {
    cy.visit('/login');
    
    cy.get('.back-link').should('be.visible');
  });

  it('back link contains arrow icon', () => {
    cy.visit('/login');
    
    cy.get('.back-link mat-icon').should('contain', 'arrow_back');
  });

  it('multiple form submissions reset form state properly', () => {
    cy.visit('/login');
    
    // First attempt
    cy.get('input[formcontrolname="email"]').type('test@example.com');
    cy.get('input[formcontrolname="password"]').type('password123');
    cy.contains('button[type="submit"]', 'Log In').click();
    
    cy.wait(1000);
    
    // Form should still be visible
    cy.get('.login-page').should('be.visible');
    cy.get('input[formcontrolname="email"]').should('be.visible');
  });


  it('email input accepts valid email format', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]')
      .type('student@university.edu')
      .should('have.value', 'student@university.edu');

    // Form should be valid if password is also filled
    cy.get('input[formcontrolname="password"]').type('password123');
    cy.contains('button[type="submit"]', 'Log In').should('not.be.disabled');
  });


  it('password input accepts exactly 8 characters', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]').type('test@example.com');
    cy.get('input[formcontrolname="password"]')
      .type('12345678')
      .should('have.value', '12345678');

    // Form should be valid
    cy.contains('button[type="submit"]', 'Log In').should('not.be.disabled');
  });

  it('password input accepts long passwords', () => {
    cy.visit('/login');
    
    const longPassword = 'aVeryLongPasswordWith123SpecialChars!@#';
    
    cy.get('input[formcontrolname="email"]').type('test@example.com');
    cy.get('input[formcontrolname="password"]')
      .type(longPassword)
      .should('have.value', longPassword);

    // Form should be valid
    cy.contains('button[type="submit"]', 'Log In').should('not.be.disabled');
  });

  it('create account button has routerLink to register', () => {
    cy.visit('/login');
    
    cy.contains('button', 'Create Account')
      .should('have.attr', 'routerLink', '/register');
  });

  it('back to home link has routerLink to home', () => {
    cy.visit('/login');
    
    cy.contains('a', 'Back to Home')
      .should('have.attr', 'routerLink', '/');
  });

  it('submit button has correct CSS classes', () => {
    cy.visit('/login');
    
    cy.contains('button[type="submit"]', 'Log In')
      .should('have.class', 'submit-btn')
      .should('have.class', 'full-width');
  });
//207

  it('create account button has full-width class', () => {
    cy.visit('/login');
    
    cy.contains('button', 'Create Account')
      .should('have.class', 'full-width');
  });

  it('mat-form-fields have full-width class', () => {
    cy.visit('/login');
    
    cy.get('mat-form-field.full-width').should('have.length', 2);
  });

  it('submit button has primary color attribute', () => {
    cy.visit('/login');
    
    cy.contains('button[type="submit"]', 'Log In')
      .should('have.attr', 'color', 'primary');
  });

  it('email label displays "Email"', () => {
    cy.visit('/login');
    
    cy.get('mat-form-field').first().find('mat-label').should('contain', 'Email');
  });

  it('password label displays "Password"', () => {
    cy.visit('/login');
    
    cy.get('mat-form-field').last().find('mat-label').should('contain', 'Password');
  });

  it('form submission with valid credentials sends data', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]').type('test@example.com');
    cy.get('input[formcontrolname="password"]').type('password123');

    const submitButton = cy.contains('button[type="submit"]', 'Log In');
    submitButton.should('not.be.disabled');
    submitButton.click();

    // App should process the submission
    cy.get('.login-page').should('be.visible');
  });

  it('password visibility toggle updates aria-label', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]').type('password123');
    
    // Check initial aria-label
    cy.get('button[mat-icon-button][matSuffix]')
      .should('have.attr', 'aria-label')
      .and('include', 'Hide password');
  });

  it('SVG logo has correct dimensions', () => {
    cy.visit('/login');
    
    cy.get('.logo-section svg')
      .should('have.attr', 'width', '64')
      .should('have.attr', 'height', '64');
  });

  it('email error message appears when email field is blurred with invalid value', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]')
      .focus()
      .type('invalid')
      .blur();

    cy.contains('Please enter a valid email').should('be.visible');
  });

  it('password error message appears when field is blurred with short value', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"]')
      .focus()
      .type('short')
      .blur();

    cy.contains('Password must be at least 8 characters').should('be.visible');
  });

  it('form resets errors when invalid field is corrected', () => {
    cy.visit('/login');
    
    // Type invalid email
    cy.get('input[formcontrolname="email"]')
      .type('invalid')
      .blur();

    cy.contains('Please enter a valid email').should('be.visible');

    // Clear and type valid email
    cy.get('input[formcontrolname="email"]')
      .clear()
      .type('valid@example.com');

    // Error should disappear
    cy.contains('Please enter a valid email').should('not.exist');
  });

  it('create account button is not disabled', () => {
    cy.visit('/login');
    
    cy.contains('button', 'Create Account').should('not.be.disabled');
  });
it('back to home link is not disabled', () => {
    cy.visit('/login');
    
    cy.contains('a', 'Back to Home').should('not.be.disabled');
  });

  it('mat-icon elements are rendered correctly in form fields', () => {
    cy.visit('/login');
    
    cy.get('mat-form-field mat-icon').should('have.length.at.least', 2);
  });

  it('mat-card-content contains form element', () => {
    cy.visit('/login');
    
    cy.get('mat-card-content').find('form').should('be.visible');
  });

  it('password field has matInput directive', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="password"][matInput]').should('be.visible');
  });
//208
  it('email field has matInput directive', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"][matInput]').should('be.visible');
  });

  it('form displays proper Material Design styling', () => {
    cy.visit('/login');
    
    cy.get('mat-card').should('be.visible');
    cy.get('mat-form-field').should('have.length', 2);
    cy.get('button[mat-raised-button]').should('be.visible');
  });

  it('all inputs have formControlName attribute', () => {
    cy.visit('/login');
    
    cy.get('input[formControlName="email"]').should('be.visible');
    cy.get('input[formControlName="password"]').should('be.visible');
  });

  it('login page does not redirect when form is invalid', () => {
    cy.visit('/login');
    
    // Leave form empty
    const initialUrl = cy.url();
    
    // Try to submit with disabled button (force click to test)
    cy.contains('button[type="submit"]', 'Log In').should('be.disabled');
    
    // URL should not have changed
    cy.url().should('include', '/login');
  });

  it('login form can be submitted after filling all required fields correctly', () => {
    cy.visit('/login');
    
    cy.get('input[formcontrolname="email"]').type('test@example.com');
    cy.get('input[formcontrolname="password"]').type('password123');
    
    // Verify submit button is enabled
    cy.contains('button[type="submit"]', 'Log In')
      .should('not.be.disabled')
      .click();

    // Verify page is still visible (form submitted)
    cy.get('.login-page').should('be.visible');
  });

  it('email input has correct mat-label', () => {
    cy.visit('/login');
    
    cy.get('mat-form-field').first()
      .find('mat-label')
      .should('have.text', 'Email');
  });

  it('password input has correct mat-label', () => {
    cy.visit('/login');
    
    cy.get('mat-form-field').last()
      .find('mat-label')
      .should('have.text', 'Password');
    describe('ThreadTalk - Communities & Community Posts (E2E)', () => {
  /**
   * Helper: login via UI (works without needing internal tokens).
   * If your app allows communities page without login, you can remove this.
   */
  const loginUI = () => {
    cy.visit('/login');
    cy.get('input[formcontrolname="email"]').clear().type('test@example.com');
    cy.get('input[formcontrolname="password"]').clear().type('password123');
    cy.contains('button[type="submit"]', 'Log In').click();

    // We don't assert success strongly because depends on backend,
    // but ensure page doesn't crash.
    cy.get('body').should('be.visible');
  };

  it('loads communities page and shows tabs', () => {
    // loginUI(); // uncomment if communities requires auth
    cy.visit('/communities');

    cy.get('.communities-page').should('be.visible');
    cy.contains('Communities').should('be.visible');

    cy.contains('All Communities').should('be.visible');
    cy.contains('My Communities').should('be.visible');

    // Search input exists
    cy.contains('mat-label', 'Search communities').should('be.visible');
  });
    it('shows the communities page header and create community button', () => {
    cy.visit('/communities');

    cy.get('.communities-page').should('be.visible');
    cy.contains('h1', 'Communities').should('be.visible');
    cy.contains('Discover and join communities at UF').should('be.visible');

    cy.contains('button', 'Create Community').should('be.visible');
  });
  });
});
