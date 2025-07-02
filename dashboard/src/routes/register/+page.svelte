<script>
  import { SignIn } from "@auth/sveltekit/components";
  import { onMount } from 'svelte';
  
  let mounted = false;
  let isLoading = false;
  let redirectTo = '/account';
  
  if (typeof window !== 'undefined') {
    const { origin } = window.location;
    redirectTo = origin + '/account';
  }
  
  onMount(() => {
    mounted = true;
  });
  
  function handleSignUp() {
    isLoading = true;
    setTimeout(() => {
      isLoading = false;
    }, 2000);
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-brand-50 via-background to-brand-100 dark:from-background dark:via-surface-elevated dark:to-surface-muted flex items-center justify-center p-4 relative overflow-hidden">
  <!-- Animated background elements -->
  <div class="absolute inset-0 overflow-hidden">
    <div class="absolute top-1/4 left-1/4 w-64 h-64 bg-gradient-to-r from-primary/10 to-accent/15 rounded-full blur-3xl animate-pulse"></div>
    <div class="absolute bottom-1/4 right-1/4 w-96 h-96 bg-gradient-to-r from-accent/10 to-brand-300/20 rounded-full blur-3xl animate-pulse" style="animation-delay: 1s;"></div>
    <div class="absolute top-1/2 right-1/3 w-48 h-48 bg-gradient-to-r from-brand-200/20 to-primary/10 rounded-full blur-3xl animate-pulse" style="animation-delay: 2s;"></div>
  </div>
  <div class="absolute inset-0 opacity-30">
    <div class="absolute inset-0 bg-gradient-to-br from-brand-100/30 via-transparent to-accent/20 dark:from-surface-elevated/20 dark:to-surface-muted/30"></div>
  </div>
  <div class="relative z-10 w-full max-w-md transform transition-all duration-700 {mounted ? 'translate-y-0 opacity-100' : 'translate-y-8 opacity-0'}">
    <div class="glass border border-border/50 rounded-3xl shadow-2xl shadow-primary/10 dark:shadow-primary/20 p-8 relative overflow-hidden">
      <div class="absolute top-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-primary/60 to-transparent animate-pulse"></div>
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-primary to-accent rounded-2xl mb-6 shadow-lg shadow-primary/25 transform hover:scale-105 transition-transform duration-300">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M16 12V8a4 4 0 00-8 0v4m8 0v4a4 4 0 01-8 0v-4" />
          </svg>
        </div>
        <h1 class="text-4xl font-bold bg-gradient-to-r from-foreground via-primary to-accent bg-clip-text text-transparent mb-3">
          Create your account
        </h1>
        <p class="text-muted-foreground text-lg font-medium">
          Sign up to get started with your dashboard
        </p>
      </div>
      <div class="space-y-6">
        <SignIn provider="hydra" options={{ prompt: "select_account", redirectTo }}>
          <button 
            slot="submitButton" 
            on:click={handleSignUp}
            disabled={isLoading}
            class="group relative w-full flex items-center justify-center gap-3 bg-gradient-to-r from-primary to-accent hover:from-primary/90 hover:to-accent/90 disabled:from-muted disabled:to-muted text-white font-semibold py-4 px-6 rounded-2xl shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/40 transform hover:scale-[1.02] active:scale-[0.98] transition-all duration-200 disabled:cursor-not-allowed disabled:transform-none overflow-hidden"
          >
            <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000"></div>
            {#if isLoading}
              <div class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              <span>Signing up...</span>
            {:else}
              <svg class="w-5 h-5" viewBox="0 0 24 24">
                <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="currentColor" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
              <span>Sign up with Google</span>
            {/if}
          </button>
        </SignIn>
        <div class="relative flex items-center my-8">
          <div class="flex-grow border-t border-border"></div>
          <span class="flex-shrink mx-4 text-muted-foreground text-sm font-medium glass px-3 py-1 rounded-full border border-border">
            Secure Authentication
          </span>
          <div class="flex-grow border-t border-border"></div>
        </div>
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
      <div class="mt-8 pt-6 border-t border-border text-center space-y-3">
        <a 
          href="/login" 
          class="inline-flex items-center text-sm text-muted-foreground hover:text-primary transition-colors duration-200 group"
        >
          <span>Already have an account?</span>
          <svg class="w-4 h-4 ml-1 transform group-hover:translate-x-1 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M9 5l7 7-7 7" />
          </svg>
        </a>
      </div>
    </div>
    <div class="absolute -top-4 -right-4 w-24 h-24 bg-gradient-to-br from-accent/20 to-brand-300/20 rounded-full blur-xl animate-pulse" style="animation-delay: 0.5s;"></div>
    <div class="absolute -bottom-4 -left-4 w-32 h-32 bg-gradient-to-br from-primary/20 to-accent/20 rounded-full blur-xl animate-pulse" style="animation-delay: 1s;"></div>
  </div>
</div>
