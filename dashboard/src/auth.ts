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
          scope: "openid offline hydra.openid introspect",
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
      // console.log("JWT Callback - Token:", token);
      // console.log("JWT Callback - Account:", account);
      // console.log("JWT Callback - User:", user);
      
      // On initial sign in, account will contain the OAuth tokens
      if (account) {
        // Store the OAuth tokens in the JWT token
        token.access_token = account.access_token;
        token.refresh_token = account.refresh_token;
        // Prefer expires_at if present, otherwise calculate from expires_in
        if (typeof account.expires_at === 'number') {
          token.access_token_expires = account.expires_at;
        } else if (typeof account.expires_in === 'number') {
          token.access_token_expires = Math.floor(Date.now() / 1000) + account.expires_in;
        } else {
          token.access_token_expires = undefined;
        }
        token.provider = account.provider;
        token.provider_account_id = account.providerAccountId;
      }
      // Check if access token has expired
      if (token.access_token_expires && Date.now() > Number(token.access_token_expires) * 1000) {
        console.log("Access token expired, need to refresh");
        return await refreshAccessToken(token);
      }
      
      return token;
    },
    
    async session({ session, token }) {
      // console.log("Session Callback - Session:", session);
      // console.log("Session Callback - Token:", token);
      
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
      // console.log("SignIn Callback - User:", user);
      // console.log("SignIn Callback - Account:", account);
      // console.log("SignIn Callback - Profile:", profile);
      
      // Return true to allow the sign in
      return true;
    }
  },
  session: {
    strategy: "jwt",
    // Optional: Set session max age (default is 30 days)
    maxAge: 30 * 24 * 60 * 60, // 30 days
  },
  debug: process.env.NODE_ENV === "development",
  trustHost: true,
});

// Add this helper function at the top or inside the jwt callback
async function refreshAccessToken(token: any) {
  try {
    const response = await fetch("http://127.0.0.1:8080/oauth2/token", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
        "Authorization": "Basic " + Buffer.from("4b41cd38-43ed-4e3a-9a88-bd384af21732:fbOoeUd9fEiw6LM~TWhg70zhTo").toString("base64"),
      },
      body: new URLSearchParams({
        grant_type: "refresh_token",
        refresh_token: token.refresh_token,
        client_id: "4b41cd38-43ed-4e3a-9a88-bd384af21732",
        client_secret: "fbOoeUd9fEiw6LM~TWhg70zhTo",
      }),
    });

    const refreshedTokens = await response.json();

    if (!response.ok) throw refreshedTokens;

    return {
      ...token,
      access_token: refreshedTokens.access_token,
      access_token_expires: Math.floor(Date.now() / 1000) + refreshedTokens.expires_in,
      refresh_token: refreshedTokens.refresh_token ?? token.refresh_token, // Some providers only return a new refresh token sometimes
    };
  } catch (error) {
    console.error("Error refreshing access token", error);
    return {
      ...token,
      error: "RefreshTokenError",
    };
  }
}