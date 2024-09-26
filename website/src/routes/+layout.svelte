<script>
  import "../app.css";
  import Navbar from "../Components/Navbar.svelte";
  import Footer from "../Components/Footer.svelte";
</script>

<div class="relative bg-gradient-to-b from-gray-50 to-gray-200 w-svw">
  <Navbar />
  <main class="w-screen">
    <slot />
  </main>

  <Footer />
</div>

<style>
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }

  main {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 1rem;
    /* width: 100%; */
    /* max-width: 64rem; */
    margin: 0 auto;
    box-sizing: border-box;
  }

  footer {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 12px;
  }

  footer a {
    font-weight: bold;
  }

  @media (min-width: 480px) {
    footer {
      padding: 12px 0;
    }
  }
</style>
