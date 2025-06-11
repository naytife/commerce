<script>
  import { SignIn } from "@auth/sveltekit/components";
  import { onMount } from 'svelte';
  
  let mounted = false;
  let isLoading = false;
  
  onMount(() => {
    mounted = true;
  });
  
  function handleSignIn() {
    isLoading = true;
    // The actual sign-in will be handled by the SignIn component
    // This is just for the loading state UI
    setTimeout(() => {
      isLoading = false;
    }, 2000);
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-brand-50 via-background to-brand-100 dark:from-background dark:via-surface-elevated dark:to-surface-muted flex items-center justify-center p-4 relative overflow-hidden">
  
  <!-- Animated background elements -->
  <div class="absolute inset-0 overflow-hidden">
    <!-- Floating orbs -->
    <div class="absolute top-1/4 left-1/4 w-64 h-64 bg-gradient-to-r from-primary/10 to-accent/15 rounded-full blur-3xl animate-pulse"></div>
    <div class="absolute bottom-1/4 right-1/4 w-96 h-96 bg-gradient-to-r from-accent/10 to-brand-300/20 rounded-full blur-3xl animate-pulse" style="animation-delay: 1s;"></div>
    <div class="absolute top-1/2 right-1/3 w-48 h-48 bg-gradient-to-r from-brand-200/20 to-primary/10 rounded-full blur-3xl animate-pulse" style="animation-delay: 2s;"></div>
  </div>

  <!-- Grid pattern overlay -->
  <div class="absolute inset-0 opacity-30">
    <div class="absolute inset-0 bg-gradient-to-br from-brand-100/30 via-transparent to-accent/20 dark:from-surface-elevated/20 dark:to-surface-muted/30"></div>
  </div>

  <!-- Main login card -->
  <div class="relative z-10 w-full max-w-md transform transition-all duration-700 {mounted ? 'translate-y-0 opacity-100' : 'translate-y-8 opacity-0'}">
    
    <!-- Glassmorphism card -->
    <div class="glass border border-border/50 rounded-3xl shadow-2xl shadow-primary/10 dark:shadow-primary/20 p-8 relative overflow-hidden">
      
      <!-- Shimmer effect on top border -->
      <div class="absolute top-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-primary/60 to-transparent animate-pulse"></div>
      
      <!-- Header section -->
      <div class="text-center mb-8">
        <!-- Logo/Icon -->
        <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-primary to-accent rounded-2xl mb-6 shadow-lg shadow-primary/25 transform hover:scale-105 transition-transform duration-300">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </div>
        
        <!-- Welcome text -->
        <h1 class="text-4xl font-bold bg-gradient-to-r from-foreground via-primary to-accent bg-clip-text text-transparent mb-3">
          Welcome Back
        </h1>
        <p class="text-muted-foreground text-lg font-medium">
          Please sign in to access your dashboard
        </p>
      </div>

      <!-- Sign in form -->
      <div class="space-y-6">
        
        <!-- Google Sign-in using your original SignIn component -->
        <SignIn provider="hydra" options={{ redirectTo: "http://localhost:5173/account" }}>
          <button 
            slot="submitButton" 
            on:click={handleSignIn}
            disabled={isLoading}
            class="group relative w-full flex items-center justify-center gap-3 bg-gradient-to-r from-primary to-accent hover:from-primary/90 hover:to-accent/90 disabled:from-muted disabled:to-muted text-white font-semibold py-4 px-6 rounded-2xl shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/40 transform hover:scale-[1.02] active:scale-[0.98] transition-all duration-200 disabled:cursor-not-allowed disabled:transform-none overflow-hidden"
          >
            <!-- Button shine effect -->
            <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000"></div>
            
            {#if isLoading}
              <div class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              <span>Signing in...</span>
            {:else}
              <svg class="w-5 h-5" viewBox="0 0 24 24">
                <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="currentColor" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
              <span>Sign in with Google</span>
            {/if}
          </button>
        </SignIn>

        <!-- Divider -->
        <div class="relative flex items-center my-8">
          <div class="flex-grow border-t border-border"></div>
          <span class="flex-shrink mx-4 text-muted-foreground text-sm font-medium glass px-3 py-1 rounded-full border border-border">
            Secure Authentication
          </span>
          <div class="flex-grow border-t border-border"></div>
        </div>

        <!-- Additional info -->
        <div class="bg-gradient-to-r from-brand-50 to-accent/10 dark:from-surface-elevated dark:to-surface-muted rounded-2xl p-4 border border-border">
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0 w-6 h-6 bg-primary rounded-full flex items-center justify-center mt-0.5">
              <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width={3} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-foreground text-sm mb-1">
                Secure & Private
              </h3>
              <p class="text-muted-foreground text-xs leading-relaxed">
                Your data is protected with enterprise-grade security. We never store your Google credentials.
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer links -->
      <div class="mt-8 pt-6 border-t border-border text-center space-y-3">
        <a 
          href="/forgot-password" 
          class="inline-flex items-center text-sm text-muted-foreground hover:text-primary transition-colors duration-200 group"
        >
          <span>Forgot Password?</span>
          <svg class="w-4 h-4 ml-1 transform group-hover:translate-x-1 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M9 5l7 7-7 7" />
          </svg>
        </a>
        
        <p class="text-sm text-muted-foreground">
          Don't have an account? 
          <a 
            href="/register" 
            class="font-semibold text-primary hover:text-accent transition-colors duration-200 underline decoration-2 underline-offset-2 hover:decoration-accent"
          >
            Create one
          </a>
        </p>
      </div>
    </div>

    <!-- Floating elements around the card -->
    <div class="absolute -top-4 -right-4 w-24 h-24 bg-gradient-to-br from-accent/20 to-brand-300/20 rounded-full blur-xl animate-pulse" style="animation-delay: 0.5s;"></div>
    <div class="absolute -bottom-4 -left-4 w-32 h-32 bg-gradient-to-br from-primary/20 to-accent/20 rounded-full blur-xl animate-pulse" style="animation-delay: 1s;"></div>
  </div>
</div>