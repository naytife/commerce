import { SvelteKitAuth } from "@auth/sveltekit";
import OryHydra from "@auth/core/providers/ory-hydra";

export const { handle, signIn, signOut } = SvelteKitAuth({
  providers: [
    OryHydra({
      id: "hydra",
      clientId: "4b41cd38-43ed-4e3a-9a88-bd384af21732",
      clientSecret: "fbOoeUd9fEiw6LM~TWhg70zhTo",
      // Fixed: Remove trailing slash to prevent double slash in well-known URL
      issuer: "http://127.0.0.1:8080",
      authorization: {
        url: "http://127.0.0.1:8080/oauth2/auth",
        params: {
          scope: "openid offline_access hydra.openid introspect",
          app_type: "dashboard",
          // Removed static state generation - let Auth.js handle it
        },
      },
      // Re-enable default checks - Auth.js needs them for proper flow
      checks: ["state", "pkce"],
    }),
  ],
  callbacks: {
    async jwt({ token, account, user }) {
      // On initial sign in, account will contain the OAuth tokens
      if (account) {
        // Store the OAuth tokens in the JWT token
        token.access_token = account.access_token;
        token.refresh_token = account.refresh_token;
        // Set access token expiry based on the actual expires_in from OAuth provider
        token.access_token_expires = Math.floor(Date.now() / 1000) + (account.expires_in || 3600);
        token.provider = account.provider;
        token.provider_account_id = account.providerAccountId;
        // Clear any stale data
        delete token.error;
        return token;
      }
      
      // Skip refresh check if we already have an error
      if (token.error) {
        return token;
      }
      
      // Check if access token needs refresh (5 minutes before expiry)
      if (token.access_token_expires && Date.now() > (Number(token.access_token_expires) - 300) * 1000) {
        const refreshedToken = await refreshAccessToken(token);
        return refreshedToken;
      }
      
      return token;
    },
    
    async session({ session, token }) {
      // Send properties to the client
      if (token) {
        // Set the access token in the session (custom property)
        (session as any).access_token = token.access_token;
        session.user = {
          id: (token.sub as string) || "",
          email: (token.provider_account_id as string) || "",
          name: (token.name as string) || undefined,
          image: (token.picture as string) || undefined,
          emailVerified: null,
        };
        // Optionally set other custom properties on the session
        (session as any).provider = token.provider;
        (session as any).provider_account_id = token.provider_account_id;
        (session as any).access_token_expires = token.access_token_expires;
      }
      
      return session;
    },

    async signIn({ user, account, profile }) {
      // Return true to allow the sign in
      return true;
    }
  },
  session: {
    strategy: "jwt",
    maxAge: 30 * 24 * 60 * 60, // 30 days
  },
  debug: false, // Disable debug in production
  trustHost: true,
});

// Token refresh management
let isRefreshing = false;
let refreshPromise: Promise<any> | null = null;
let lastRefreshTime = 0;
let cachedRefreshedToken: any = null;
let lastTokenId: string | null = null;

async function refreshAccessToken(token: any) {
  const now = Date.now() / 1000;
  const tokenId = token.jti || token.sub || 'unknown';
  
  // Check for stale token first - if this token ID is older than our last refresh
  if (lastTokenId && tokenId !== lastTokenId && cachedRefreshedToken) {
    return cachedRefreshedToken;
  }
  
  // If we have a cached refreshed token that's still valid, use it
  if (cachedRefreshedToken && 
      cachedRefreshedToken.access_token_expires && 
      now < (Number(cachedRefreshedToken.access_token_expires) - 300)) { // 5 minutes buffer
    return cachedRefreshedToken;
  }
  
  // If refresh is already in progress, wait for it
  if (isRefreshing && refreshPromise) {
    try {
      const result = await refreshPromise;
      return result;
    } catch (error) {
      // If the shared refresh failed, return error token
      return {
        ...token,
        error: "RefreshTokenError",
      };
    }
  }
  
  // Prevent too frequent refresh attempts (minimum 30 seconds between attempts)
  if (now - lastRefreshTime < 30) {
    return token;
  }
  
  try {
    isRefreshing = true;
    lastRefreshTime = now;
    lastTokenId = tokenId;
    
    // Create a promise that other concurrent calls can wait for
    refreshPromise = performTokenRefresh(token);
    const result = await refreshPromise;
    
    // Cache the result for immediate use by other concurrent calls
    cachedRefreshedToken = result;
    return result;
    
  } catch (error) {
    // If refresh failed, return error token to trigger logout
    console.error("Token refresh failed:", error);
    return {
      ...token,
      error: "RefreshTokenError",
    };
  } finally {
    isRefreshing = false;
    // Don't clear refreshPromise immediately - let other waiters get the result
    setTimeout(() => {
      refreshPromise = null;
    }, 100);
  }
}

async function performTokenRefresh(token: any) {
  if (!token.refresh_token) {
    console.error("No refresh token available");
    return {
      ...token,
      error: "RefreshTokenError",
    };
  }

  const response = await fetch("http://127.0.0.1:8080/oauth2/token", {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
      "Authorization": "Basic " + Buffer.from("4b41cd38-43ed-4e3a-9a88-bd384af21732:fbOoeUd9fEiw6LM~TWhg70zhTo").toString("base64"),
    },
    body: new URLSearchParams({
      grant_type: "refresh_token",
      refresh_token: token.refresh_token,
    }),
  });

  const refreshedTokens = await response.json();

  if (!response.ok) {
    console.error("Failed to refresh token:", refreshedTokens.error || "Unknown error");
    throw refreshedTokens;
  }

  return {
    ...token,
    access_token: refreshedTokens.access_token,
    access_token_expires: Math.floor(Date.now() / 1000) + refreshedTokens.expires_in,
    refresh_token: refreshedTokens.refresh_token ?? token.refresh_token,
    error: undefined,
  };
}