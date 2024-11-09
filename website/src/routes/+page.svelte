<script>
  import Seperator from "../Components/Seperator.svelte";
  import Dashboardscreen1 from "$lib/images/dashboard-svg1.svg";
  import Dashboardscreen2 from "$lib/images/side-530Asset7.svg";
  import Dashboardstarwhite from "$lib/images/side-530Assetscreen2.svg";
  import Dashboardstarpurple from "$lib/images/side-530Asset11.svg";
  import Dashboardstarred from "$lib/images/side-530Asset8.svg";

  import gif from "$lib/images/SVG/Cloudstorage.gif";
  import Pricing from "../Components/Pricing.svelte";

  // Array of animations with 16 problems and solutions related to creating an online store
  const animations = [
    {
      problem: "High cart abandonment rates",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/shopping-cart.svg", // Shopping cart outline icon
      solution: "Optimize the checkout process and use retargeting",
    },
    {
      problem: "Slow website loading times",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/clock.svg", // Slow clock outline icon
      solution: "Use a CDN and optimize images and scripts",
    },
    {
      problem: "Poor mobile optimization",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/device-mobile.svg", // Mobile phone outline icon
      solution: "Use responsive design and mobile-friendly navigation",
    },
    {
      problem: "Security vulnerabilities",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/lock-closed.svg", // Lock outline icon
      solution: "Implement SSL and regular security audits",
    },
    {
      problem: "Difficulties managing inventory",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/archive.svg", // Inventory box outline icon
      solution: "Use automated inventory management tools",
    },
    {
      problem: "High shipping costs",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/truck.svg", // Truck outline icon
      solution: "Offer free shipping or bundle shipping costs",
    },
    {
      problem: "Poor SEO performance",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/chart-bar.svg", // Chart bar outline icon (SEO performance)
      solution: "Optimize product pages and blog posts for search engines",
    },
    {
      problem: "Low conversion rates",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/trending-down.svg", // Conversion rates outline icon
      solution: "Use A/B testing and improve user experience",
    },
    {
      problem: "Handling customer complaints",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/support.svg", // Customer support outline icon
      solution: "Improve customer support and communication channels",
    },
    {
      problem: "Complicated returns process",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/switch-horizontal.svg", // Returns process outline icon
      solution: "Simplify return policies and offer easy return options",
    },
    {
      problem: "Limited payment options",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/credit-card.svg", // Credit card outline icon (payment options)
      solution: "Offer multiple payment gateways like PayPal and credit cards",
    },
    {
      problem: "Lack of customer trust",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/shield-check.svg", // Trust shield outline icon
      solution: "Show customer reviews and offer guarantees",
    },
    {
      problem: "Difficulty scaling the business",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/scale.svg", // Scale business outline icon
      solution: "Automate processes and outsource non-core tasks",
    },
    {
      problem: "High bounce rates",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/arrow-up.svg", // Bounce rates outline icon
      solution: "Improve landing pages and content relevance",
    },
    {
      problem: "Product pages lacking information",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/information-circle.svg", // Information outline icon
      solution: "Add detailed product descriptions and high-quality images",
    },
    {
      problem: "Inconsistent branding",
      icon: "https://cdn.jsdelivr.net/npm/heroicons/outline/tag.svg", // Branding outline icon
      solution: "Maintain consistent branding across all channels",
    },
  ];

  const getRandomLightColor = () => {
    // Generate a random light color with alpha transparency
    const r = Math.floor(Math.random() * 256);
    const g = Math.floor(Math.random() * 256);
    const b = Math.floor(Math.random() * 256);
    const alpha = 0.2; // Adjust alpha value (0.0 to 1.0) for transparency
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
  };

  // Generate random delays for animations
  const generateRandomDelay = () => Math.random() * (10 - 5) + 5; // Generates a delay between 5 and 10 seconds
</script>

<svelte:head>
  <title>Home</title>
  <meta name="description" content="Naytife app" />
</svelte:head>

<section>
  <div class="flex items-center text-center space-y-5 lg:space-y- flex-col">
    <div class="lg:w-2/3 relative">
      <!-- <div
				class="w-48 rounded-full top-2/4 transform translate-y-[100%] translate-x-[223%] absolute h-5 ring-teal-400 ring-[10px]"
			></div> -->
      <div
        class="w-52 lg:flex hidden rounded-full top-2/4 transform -translate-y-[265%] translate-x-[140%] absolute h-5 ring-orange-300 ring-[10px]"
      ></div>
      <div
        class="w-64 lg:flex hidden rounded-full top-2/4 transform translate-y-[100%] translate-x-[15%] absolute h-5 ring-purple-300 ring-[10px]"
      ></div>
      <class class="text-3xl lg:text-5xl flex-row font-serif relative">
        <p class="font-extrabold text-gray-700">
          Grow your business more efficiently with <i class="text-gray-600"
            >Naytife</i
          >
        </p>
      </class>

      <h2 class="text-black font-serif text-sm text-wrap">
        Enjoy reliable, scalable, fast and secure, Hosting for your business
        --managed for you
      </h2>
    </div>
    <div class="relative flex items-center">
      <div>
        <img
          class="rounded-3xl"
          src="https://images.rawpixel.com/image_800/cHJpdmF0ZS9sci9pbWFnZXMvd2Vic2l0ZS8yMDI0LTA1L2FyLTM5Mi0wODA1MjQtYWRqLTAwMi1sdzBmY2sxcC5qcGc.jpg"
          alt="dashboardscreen1"
        />
      </div>
      <img
        class="absolute hidden lg:flex w-[13rem] top-1/3 left-0 transform drop-shadow-2xl animate-custom-fade-1 -translate-x-1/2 -translate-y-1/2"
        src={Dashboardscreen1}
        alt="dashboardscreen1"
      />
      <img
        class="absolute hidden lg:flex w-[13rem] top-2/3 drop-shadow-2xl right-0 transform animate-custom-fade-1 translate-x-1/2 -translate-y-1/2"
        src={Dashboardscreen2}
        alt="dashboardscreen2"
      />
      <img
        class="absolute w-[1.5rem] top-1/3 right-0 transform translate-x-1/4 animate-pop"
        src={Dashboardstarwhite}
        alt="dashboardstarwhite"
      />
      <img
        class="absolute w-[2.3rem] top-[20%] right-0 hidden lg:flex transform translate-x-[400%] animate-custom-fade-2"
        src={Dashboardstarred}
        alt="dashboardstarred"
      />
      <img
        class="absolute w-[2rem] top-2/3 left-0 transform -translate-x-[300%] hidden lg:flex animate-custom-fade-5"
        src={Dashboardstarpurple}
        alt="dashboardstarpurple"
      />
      <img
        class="absolute w-[1.5rem] top-[80%] left-0 transform -translate-x-[20%] animate-pop"
        src={Dashboardstarwhite}
        alt="dashboardstarwhite"
      />
    </div>
  </div>
</section>

<section class="flex flex-col items-center justify-center space-y-12">
  <div class="relative md:max-w-[768px] mx-auto pt-9 pb-6">
    <div class="flex justify-center">
      <div
        class="flex items-center justify-center lg:w-1/5 rounded-full gap-2 bg-gray-50 py-1 px-3 border border-gray-200 shadow-md"
      >
        <img class="w-5 h-5" src={gif} alt="" />
        <button class="text-xs flex">The problem</button>
      </div>
    </div>
    <div class="text-center">
      <p class="text-4xl lg:text-6xl my-4 font-bold text-gray-700">
        The internet is complicated
      </p>
    </div>
  </div>

  <!-- Animated Items Section -->
  <div
    class="relative w-full flex items-center justify-center overflow-x-auto lg:w-max"
  >
    <div class="w-max lg:overflow-hidden">
      <div class="grid grid-cols-2 lg:grid-cols-8 gap-2 md:gap-6 relative">
        <!-- Left and Right Gradient Effects (Hidden on Mobile) -->
        <div
          class="hidden lg:block bg-gradient-to-r z-10 from-gray-100/10 to-gray-100 absolute right-0 w-1/6 h-full"
        ></div>
        <div
          class="hidden lg:block bg-gradient-to-l z-10 from-gray-100/10 to-gray-100 absolute left-0 w-1/6 h-full"
        ></div>

        <!-- Each item with animations and colors -->
        {#each animations as item, i}
          <div
            style="--animation-delay: {generateRandomDelay()}s; --bg-color: {getRandomLightColor()};"
            class="lg:grid-item flex items-center justify-evenly pt-5 flex-col text-center border border-gray-100 rounded-xl problem-item w-40 h-40"
          >
            <!-- Icon animation with random background color -->
            <div
              style="background-color: var(--bg-color);"
              class="rounded-full problem-item-icon object-contain w-2/4 flex justify-center items-center h-1/3"
            >
              <img
                src={item.icon}
                alt={item.problem}
                class="problem-item-icon-img object-contain"
              />
            </div>

            <!-- Label animation -->
            <div class="text-sm font-medium problem-item-label">
              {item.problem}
            </div>
          </div>
        {/each}
      </div>
    </div>
  </div>

  <Seperator />
</section>

<section class="flex flex-col items-center justify-center space-y-8">
  <div class="relative text-center mx-auto pt-9 pb-6">
    <div class="flex justify-center">
      <div
        class="flex items-center justify-center w-1/3 lg:w-1/5 rounded-full gap-2 bg-gray-100 py-1 px-3 border border-gray-200 shadow-md"
      >
        <img class="w-5 h-5" src={gif} alt="" />
        <button class="text-xs text-black">The solution</button>
      </div>
    </div>
    <div class=" lg:px-56">
      <p class="text-4xl lg:text-6xl my-4 font-bold text-gray-700">
        Naytife is a complete user control platform for OSM
      </p>
    </div>
    <h2 class="text-sm lg:text-base text-gray-700">
      Take a look at how Naytife can save you time and money
    </h2>
  </div>

  <!-- Get Started Button -->
  <div class="flex justify-center">
    <button
      class="bg-green-500 text-white text-sm lg:text-base py-2 px-6 rounded-full hover:bg-green-600"
    >
      Get Started
    </button>
  </div>

  <Seperator />
</section>

<section class="flex space-y-5">
  <div
    class="relative savings-optimization items-center text-center flex flex-col"
  >
    <div
      class="relative rounded-lg w-full lg:w-[calc(100%-40%)] h-full mb-10 overflow-hidden"
    >
      <!-- Background Canvas -->
      <canvas class="dotted-background canvas-small-screen h-full w-full"
      ></canvas>

      <div
        class="absolute top-0 md:left-1/2 md:top-8 left-[calc(100%-55%)] right-1/2 transform -translate-x-1/2 grid grid-cols-3 grid-rows-4 h-[calc(100%-30%)] md:h-[calc(100%-35%)] w-[calc(100%)] md:w-[70%]"
      >
        <!-- Logo  -->
        <div
          class="rounded-2xl flex w-28 h-full row-start-1 items-center justify-center row-end-4 row-span-2 col-start-1"
        >
          <div class="h-full space-y-20 w-full flex flex-col">
            <div class="logo-item">
              <div class="relative">
                <div
                  class="savings-optimization-logos-item absolute z-10"
                  data-step="3"
                >
                  <div
                    class="-inset-x-4 -inset-y-2 absolute bg-transparent rounded-3xl -z-10"
                  ></div>
                  <div
                    class="mb-2 mx-auto w-max relative savings-optimization-logos-item-icon transition-all p-5 rounded-full bg-white"
                  >
                    <div
                      class="overlay savings-optimization-logos-item-border-wrapper duration-[600ms] transition-all scale-[.8] opacity-0 grayscale"
                    >
                      <div
                        class="savings-optimization-logos-item-border bg-green-500/5 absolute -inset-8 rounded-full"
                      ></div>
                      <div class="absolute blur-[3px] -inset-8 rounded-full">
                        <div
                          class="savings-optimization-logos-item-border rounded-full w-full h-full"
                        ></div>
                      </div>
                    </div>
                    <div
                      class="savings-optimization-logos-item-check absolute duration-[600ms] -top-2 -right-2 translate-y-6 grayscale blur-[1px] opacity-0 w-5 h-5 p-1 flex items-center transition-all justify-center rounded-full"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 16 16"
                        fill="none"
                      >
                        <path
                          d="M1 8l4 4L15 4"
                          stroke="white"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        />
                      </svg>
                    </div>
                    <svg
                      class="w-8 h-8 text-green-500"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path d="M9 12h6M9 8h6M9 16h6"></path>
                    </svg>
                  </div>
                  <div class="text-xs text-center w-full text-black">
                    Product Recs
                  </div>
                  <div
                    class="savings-optimization-logos-item-optimized translate-y-12 scale-[.9] opacity-0 grayscale w-10 transition-all duration-[600ms] rounded-full center-x -top-14 flex items-center justify-center h-5"
                  >
                    <div class="p-1">
                      <div
                        class="savings-optimization-logos-item-optimized-check translate-x-9 w-8 h-8 flex items-center transition-all justify-center rounded-full duration-[600ms]"
                      >
                        <svg
                          fill="none"
                          height="8"
                          viewBox="0 0 8 8"
                          width="8"
                          xmlns="http://www.w3.org/2000/svg"
                          ><path
                            d="M1.5 4.35714L3.6 6.5L8.5 1.5"
                            stroke="white"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="1.5"
                            style="stroke: white; stroke-opacity: 1;"
                          ></path></svg
                        >
                      </div>
                    </div>
                    <div
                      class="savings-optimization-logos-item-optimized-label duration-[600ms] transition-all opacity-0 text-sm -translate-x-4 pr-6 pl-4"
                    >
                      Optimized
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="logo-item self-end">
              <div class="relative self-center">
                <div
                  class="savings-optimization-logos-item absolute z-10"
                  data-step="3"
                >
                  <div
                    class="-inset-x-4 -inset-y-2 absolute bg-transparent rounded-3xl -z-10"
                  ></div>
                  <div
                    class="mb-2 mx-auto w-max relative savings-optimization-logos-item-icon transition-all p-5 rounded-full bg-white"
                  >
                    <div
                      class="overlay savings-optimization-logos-item-border-wrapper duration-[600ms] transition-all scale-[.8] opacity-0 grayscale"
                    >
                      <div
                        class="savings-optimization-logos-item-border bg-green-500/5 absolute -inset-8 rounded-full"
                      ></div>
                      <div class="absolute blur-[3px] -inset-8 rounded-full">
                        <div
                          class="savings-optimization-logos-item-border rounded-full w-full h-full"
                        ></div>
                      </div>
                    </div>
                    <div
                      class="savings-optimization-logos-item-check absolute duration-[600ms] -top-2 -right-2 translate-y-6 grayscale blur-[1px] opacity-0 w-5 h-5 p-1 flex items-center transition-all justify-center rounded-full"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 16 16"
                        fill="none"
                      >
                        <path
                          d="M1 8l4 4L15 4"
                          stroke="white"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        />
                      </svg>
                    </div>
                    <svg
                      class="w-8 h-8 text-blue-500"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path d="M3 3v18h18V3H3zm7 7V6h4v4h3l-5 5-5-5h3z"></path>
                    </svg>
                  </div>
                  <div class="text-xs text-center text-black">
                    Scalable Hosting
                  </div>
                  <div
                    class="savings-optimization-logos-item-optimized translate-y-12 scale-[.9] opacity-0 grayscale w-10 transition-all duration-[600ms] rounded-full center-x -top-14 flex items-center justify-center h-5"
                  >
                    <div class="p-1">
                      <div
                        class="savings-optimization-logos-item-optimized-check translate-x-9 w-8 h-8 flex items-center transition-all justify-center rounded-full duration-[600ms]"
                      >
                        <svg
                          fill="none"
                          height="8"
                          viewBox="0 0 8 8"
                          width="8"
                          xmlns="http://www.w3.org/2000/svg"
                          ><path
                            d="M1.5 4.35714L3.6 6.5L8.5 1.5"
                            stroke="white"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="1.5"
                            style="stroke: white; stroke-opacity: 1;"
                          ></path></svg
                        >
                      </div>
                    </div>
                    <div
                      class="savings-optimization-logos-item-optimized-label duration-[600ms] transition-all opacity-0 text-sm -translate-x-4 pr-6 pl-4"
                    >
                      Optimized
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div
          class="rounded-2xl flex w-28 h-full justify-center items-center row-start-2 row-end-5 row-span-2 col-start-2"
        >
          <div class="h-full space-y-20 justify-between w-full flex flex-col">
            <div class="logo-item self-center">
              <div class="relative self-center">
                <div
                  class="savings-optimization-logos-item absolute z-10"
                  data-step="3"
                >
                  <div
                    class="-inset-x-4 -inset-y-2 absolute bg-transparent rounded-3xl -z-10"
                  ></div>
                  <div
                    class="mb-2 mx-auto w-max relative savings-optimization-logos-item-icon transition-all p-7 rounded-xl bg-white"
                  >
                    <div
                      class="overlay savings-optimization-logos-item-border-wrapper duration-[600ms] transition-all scale-[.8] opacity-0 grayscale"
                    >
                      <div
                        class="savings-optimization-logos-item-border bg-green-500/5 absolute -inset-8 rounded-full"
                      ></div>
                      <div class="absolute blur-[3px] -inset-8 rounded-full">
                        <div
                          class="savings-optimization-logos-item-border rounded-full w-full h-full"
                        ></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="logo-item self-center">
              <div class="relative self-center">
                <div
                  class="savings-optimization-logos-item absolute z-10"
                  data-step="3"
                >
                  <div
                    class="-inset-x-4 -inset-y-2 absolute bg-transparent rounded-3xl -z-10"
                  ></div>
                  <div
                    class="mb-2 mx-auto w-max relative savings-optimization-logos-item-icon transition-all p-5 rounded-full bg-white"
                  >
                    <div
                      class="overlay savings-optimization-logos-item-border-wrapper duration-[600ms] transition-all scale-[.8] opacity-0 grayscale"
                    >
                      <div
                        class="savings-optimization-logos-item-border bg-green-500/5 absolute -inset-8 rounded-full"
                      ></div>
                      <div class="absolute blur-[3px] -inset-8 rounded-full">
                        <div
                          class="savings-optimization-logos-item-border rounded-full w-full h-full"
                        ></div>
                      </div>
                    </div>
                    <div
                      class="savings-optimization-logos-item-check absolute duration-[600ms] -top-2 -right-2 translate-y-6 grayscale blur-[1px] opacity-0 w-5 h-5 p-1 flex items-center transition-all justify-center rounded-full"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 16 16"
                        fill="none"
                      >
                        <path
                          d="M1 8l4 4L15 4"
                          stroke="white"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        />
                      </svg>
                    </div>
                    <svg
                      class="w-8 h-8 text-yellow-500"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path d="M16.5 7.5h4M14 5v4M5 7l5-5 5 5M9 16v4h6v-4"
                      ></path>
                    </svg>
                  </div>
                  <div class="text-xs text-center text-black">
                    SEO Optimization
                  </div>
                  <div
                    class="savings-optimization-logos-item-optimized translate-y-12 scale-[.9] opacity-0 grayscale w-10 transition-all duration-[600ms] rounded-full center-x -top-14 flex items-center justify-center h-5"
                  >
                    <div class="p-1">
                      <div
                        class="savings-optimization-logos-item-optimized-check translate-x-9 w-8 h-8 flex items-center transition-all justify-center rounded-full duration-[600ms]"
                      >
                        <svg
                          fill="none"
                          height="8"
                          viewBox="0 0 8 8"
                          width="8"
                          xmlns="http://www.w3.org/2000/svg"
                          ><path
                            d="M1.5 4.35714L3.6 6.5L8.5 1.5"
                            stroke="white"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="1.5"
                            style="stroke: white; stroke-opacity: 1;"
                          ></path></svg
                        >
                      </div>
                    </div>
                    <div
                      class="savings-optimization-logos-item-optimized-label duration-[600ms] transition-all opacity-0 text-sm -translate-x-4 pr-6 pl-4"
                    >
                      Optimized
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div
          class="rounded-2xl flex w-28 h-full row-start-1 items-center justify-center row-end-4 row-span-2 col-start-3"
        >
          <div class="h-full space-y-20 w-full flex flex-col">
            <div class="logo-item self-end">
              <div class="logo-item self-end">
                <div class="relative">
                  <div
                    class="savings-optimization-logos-item absolute z-10"
                    data-step="3"
                  >
                    <div
                      class="-inset-x-4 -inset-y-2 absolute bg-transparent rounded-3xl -z-10"
                    ></div>
                    <div
                      class="mb-2 mx-auto w-max relative savings-optimization-logos-item-icon transition-all p-5 rounded-full bg-white"
                    >
                      <div
                        class="overlay savings-optimization-logos-item-border-wrapper duration-[600ms] transition-all scale-[.8] opacity-0 grayscale"
                      >
                        <div
                          class="savings-optimization-logos-item-border bg-green-500/5 absolute -inset-8 rounded-full"
                        ></div>
                        <div class="absolute blur-[3px] -inset-8 rounded-full">
                          <div
                            class="savings-optimization-logos-item-border rounded-full w-full h-full"
                          ></div>
                        </div>
                      </div>
                      <div
                        class="savings-optimization-logos-item-check absolute duration-[600ms] -top-2 -right-2 translate-y-6 grayscale blur-[1px] opacity-0 w-5 h-5 p-1 flex items-center transition-all justify-center rounded-full"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="16"
                          height="16"
                          viewBox="0 0 16 16"
                          fill="none"
                        >
                          <path
                            d="M1 8l4 4L15 4"
                            stroke="white"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                          />
                        </svg>
                      </div>
                      <svg
                        class="w-8 h-8 text-purple-500"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path d="M8 21h4a8 8 0 10-4 0zm3-16h1a5 5 0 00-5 5h3"
                        ></path>
                      </svg>
                    </div>
                    <div class="text-xs text-center text-black">
                      Customer Analytics
                    </div>
                    <div
                      class="savings-optimization-logos-item-optimized translate-y-12 scale-[.9] opacity-0 grayscale w-10 transition-all duration-[600ms] rounded-full center-x -top-14 flex items-center justify-center h-5"
                    >
                      <div class="p-1">
                        <div
                          class="savings-optimization-logos-item-optimized-check translate-x-9 w-8 h-8 flex items-center transition-all justify-center rounded-full duration-[600ms]"
                        >
                          <svg
                            fill="none"
                            height="8"
                            viewBox="0 0 8 8"
                            width="8"
                            xmlns="http://www.w3.org/2000/svg"
                            ><path
                              d="M1.5 4.35714L3.6 6.5L8.5 1.5"
                              stroke="white"
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="1.5"
                              style="stroke: white; stroke-opacity: 1;"
                            ></path></svg
                          >
                        </div>
                      </div>
                      <div
                        class="savings-optimization-logos-item-optimized-label duration-[600ms] transition-all opacity-0 text-sm -translate-x-4 pr-6 pl-4"
                      >
                        Optimized
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="logo-item self-start">
              <div class="relative self-center">
                <div
                  class="savings-optimization-logos-item absolute z-10"
                  data-step="3"
                >
                  <div
                    class="-inset-x-4 -inset-y-2 absolute bg-transparent rounded-3xl -z-10"
                  ></div>
                  <div
                    class="mb-2 mx-auto w-max relative savings-optimization-logos-item-icon transition-all p-5 rounded-full bg-white"
                  >
                    <div
                      class="overlay savings-optimization-logos-item-border-wrapper duration-[600ms] transition-all scale-[.8] opacity-0 grayscale"
                    >
                      <div
                        class="savings-optimization-logos-item-border bg-green-500/5 absolute -inset-8 rounded-full"
                      ></div>
                      <div class="absolute blur-[3px] -inset-8 rounded-full">
                        <div
                          class="savings-optimization-logos-item-border rounded-full w-full h-full"
                        ></div>
                      </div>
                    </div>
                    <div
                      class="savings-optimization-logos-item-check absolute duration-[600ms] -top-2 -right-2 translate-y-6 grayscale blur-[1px] opacity-0 w-5 h-5 p-1 flex items-center transition-all justify-center rounded-full"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 16 16"
                        fill="none"
                      >
                        <path
                          d="M1 8l4 4L15 4"
                          stroke="white"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        />
                      </svg>
                    </div>
                  </div>
                  <div class="text-xs text-center text-black">EC2</div>
                  <div
                    class="savings-optimization-logos-item-optimized translate-y-12 scale-[.9] opacity-0 grayscale w-10 transition-all duration-[600ms] rounded-full center-x -top-14 flex items-center justify-center h-5"
                  >
                    <div class="p-1">
                      <div
                        class="savings-optimization-logos-item-optimized-check translate-x-9 w-8 h-8 flex items-center transition-all justify-center rounded-full duration-[600ms]"
                      >
                        <svg
                          fill="none"
                          height="8"
                          viewBox="0 0 8 8"
                          width="8"
                          xmlns="http://www.w3.org/2000/svg"
                          ><path
                            d="M1.5 4.35714L3.6 6.5L8.5 1.5"
                            stroke="white"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="1.5"
                            style="stroke: white; stroke-opacity: 1;"
                          ></path></svg
                        >
                      </div>
                    </div>
                    <div
                      class="savings-optimization-logos-item-optimized-label duration-[600ms] transition-all opacity-0 text-sm -translate-x-4 pr-6 pl-4"
                    >
                      Optimized
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- Add more logos here as needed -->
      </div>
    </div>
    <div class="text-center px-10 md:px-40 md:w-1/2 w-full pb-6">
      <h2 class="mb-6 text-black font-bold text-xl">
        AI-Powered Online Store Optimization
      </h2>
      <p class="text-body-medium-regular text-neutral-alpha-mute-alpha">
        Use AI to streamline your store's creation and hosting with advanced
        tools like product recommendations, scalable hosting, SEO, and
        analytics.
      </p>
    </div>
  </div>

  <div
    class="text-center pt-6 mx-auto relative before:!inset-y-7 md:before:!inset-y-10 before:inside-border-x before:border-gray-300"
  >
    <div class="relative">
      <div
        class="flex items-center justify-center mb-24 h-40 relative sm:before:hidden before:inside-border-y before:border-gray-300 sm:rounded-xl"
      >
        <!-- Main Content Area -->
        <div
          class="relative max-x-0 before:hidden before:inside-border before:border-gray-300 rounded-xl"
        >
          <!-- Main Content Inner Container -->
          <div
            class="w-full gap-4 h-full z-50 flex justify-center items-center mx-auto"
          >
            <!-- Hidden Space for Small Screens -->
            <div class="md:hidden w-[18px]"></div>

            <!-- Image 1 -->
            <div
              class="w-16 h-16 flex justify-center items-center rounded-full bg-white"
            >
              <svg
                class="w-10 h-10"
                fill="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12z"
                  clip-rule="evenodd"
                ></path>
              </svg>
            </div>

            <!-- Image 2 -->
            <div
              class="w-16 h-16 flex justify-center items-center rounded-full bg-white"
            >
              <svg
                class="w-10 h-10"
                fill="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M12.315 2c2.43 0 2.784.013 3.808.06 1.064.049 1.791.218 2.427.465a4.902 4.902 0 011.772 1.153 4.902 4.902 0 011.153 1.772c.247.636.416 1.363.465 2.427.048 1.067.06 1.407.06 4.123v.08c0 2.643-.012 2.987-.06 4.043-.049 1.064-.218 1.791-.465 2.427a4.902 4.902 0 01-1.153 1.772 4.902 4.902 0 01-1.772 1.153c-.636.247-1.363.416-2.427.465-1.067.048-1.407.06-4.123.06h-.08c-2.643 0-2.987-.012-4.043-.06-1.064-.049-1.791-.218-2.427-.465a4.902 4.902 0 01-1.772-1.153 4.902 4.902 0 01-1.153-1.772c-.247-.636-.416-1.363-.465-2.427-.047-1.024-.06-1.379-.06-3.808v-.63c0-2.43.013-2.784.06-3.808.049-1.064.218-1.791.465-2.427a4.902 4.902 0 011.153-1.772A4.902 4.902 0 015.45 2.525c.636-.247 1.363-.416 2.427-.465C8.901 2.013 9.256 2 11.685 2h.63zm-.081 1.802h-.468c-2.456 0-2.784.011-3.807.058-.975.045-1.504.207-1.857.344-.467.182-.8.398-1.15.748-.35.35-.566.683-.748 1.15-.137.353-.3.882-.344 1.857-.047 1.023-.058 1.351-.058 3.807v.468c0 2.456.011 2.784.058 3.807.045.975.207 1.504.344 1.857.182.466.399.8.748 1.15.35.35.683.566 1.15.748.353.137.882.3 1.857.344 1.054.048 1.37.058 4.041.058h.08c2.597 0 2.917-.01 3.96-.058.976-.045 1.505-.207 1.858-.344.466-.182.8-.398 1.15-.748.35-.35.566-.683.748-1.15.137-.353.3-.882.344-1.857.048-1.055.058-1.37.058-4.041v-.08c0-2.597-.01-2.917-.058-3.96-.045-.976-.207-1.505-.344-1.858a3.097 3.097 0 00-.748-1.15 3.098 3.098 0 00-1.15-.748c-.353-.137-.882-.3-1.857-.344-1.023-.047-1.351-.058-3.807-.058zM12 6.865a5.135 5.135 0 110 10.27 5.135 5.135 0 010-10.27zm0 1.802a3.333 3.333 0 100 6.666 3.333 3.333 0 000-6.666zm5.338-3.205a1.2 1.2 0 110 2.4 1.2 1.2 0 010-2.4z"
                  clip-rule="evenodd"
                ></path>
              </svg>
            </div>

            <!-- Animated Antimetal Element -->
            <div
              class="savings-risk-animation-antimetal w-18 h-18 relative bg-white/64 backdrop-blur-[6px] rounded-xl"
            ></div>

            <!-- Image 3 -->
            <div
              class="w-16 h-16 flex justify-center items-center rounded-full bg-white"
            >
              <svg
                class="w-10 h-10"
                fill="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M12.315 2c2.43 0 2.784.013 3.808.06 1.064.049 1.791.218 2.427.465a4.902 4.902 0 011.772 1.153 4.902 4.902 0 011.153 1.772c.247.636.416 1.363.465 2.427.048 1.067.06 1.407.06 4.123v.08c0 2.643-.012 2.987-.06 4.043-.049 1.064-.218 1.791-.465 2.427a4.902 4.902 0 01-1.153 1.772 4.902 4.902 0 01-1.772 1.153c-.636.247-1.363.416-2.427.465-1.067.048-1.407.06-4.123.06h-.08c-2.643 0-2.987-.012-4.043-.06-1.064-.049-1.791-.218-2.427-.465a4.902 4.902 0 01-1.772-1.153 4.902 4.902 0 01-1.153-1.772c-.247-.636-.416-1.363-.465-2.427-.047-1.024-.06-1.379-.06-3.808v-.63c0-2.43.013-2.784.06-3.808.049-1.064.218-1.791.465-2.427a4.902 4.902 0 011.153-1.772A4.902 4.902 0 015.45 2.525c.636-.247 1.363-.416 2.427-.465C8.901 2.013 9.256 2 11.685 2h.63zm-.081 1.802h-.468c-2.456 0-2.784.011-3.807.058-.975.045-1.504.207-1.857.344-.467.182-.8.398-1.15.748-.35.35-.566.683-.748 1.15-.137.353-.3.882-.344 1.857-.047 1.023-.058 1.351-.058 3.807v.468c0 2.456.011 2.784.058 3.807.045.975.207 1.504.344 1.857.182.466.399.8.748 1.15.35.35.683.566 1.15.748.353.137.882.3 1.857.344 1.054.048 1.37.058 4.041.058h.08c2.597 0 2.917-.01 3.96-.058.976-.045 1.505-.207 1.858-.344.466-.182.8-.398 1.15-.748.35-.35.566-.683.748-1.15.137-.353.3-.882.344-1.857.048-1.055.058-1.37.058-4.041v-.08c0-2.597-.01-2.917-.058-3.96-.045-.976-.207-1.505-.344-1.858a3.097 3.097 0 00-.748-1.15 3.098 3.098 0 00-1.15-.748c-.353-.137-.882-.3-1.857-.344-1.023-.047-1.351-.058-3.807-.058zM12 6.865a5.135 5.135 0 110 10.27 5.135 5.135 0 010-10.27zm0 1.802a3.333 3.333 0 100 6.666 3.333 3.333 0 000-6.666zm5.338-3.205a1.2 1.2 0 110 2.4 1.2 1.2 0 010-2.4z"
                  clip-rule="evenodd"
                ></path>
              </svg>
            </div>

            <!-- Image 4 -->
            <div
              class="w-16 h-16 flex justify-center items-center rounded-full bg-white"
            >
              <svg
                class="w-10 h-10"
                fill="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M12.315 2c2.43 0 2.784.013 3.808.06 1.064.049 1.791.218 2.427.465a4.902 4.902 0 011.772 1.153 4.902 4.902 0 011.153 1.772c.247.636.416 1.363.465 2.427.048 1.067.06 1.407.06 4.123v.08c0 2.643-.012 2.987-.06 4.043-.049 1.064-.218 1.791-.465 2.427a4.902 4.902 0 01-1.153 1.772 4.902 4.902 0 01-1.772 1.153c-.636.247-1.363.416-2.427.465-1.067.048-1.407.06-4.123.06h-.08c-2.643 0-2.987-.012-4.043-.06-1.064-.049-1.791-.218-2.427-.465a4.902 4.902 0 01-1.772-1.153 4.902 4.902 0 01-1.153-1.772c-.247-.636-.416-1.363-.465-2.427-.047-1.024-.06-1.379-.06-3.808v-.63c0-2.43.013-2.784.06-3.808.049-1.064.218-1.791.465-2.427a4.902 4.902 0 011.153-1.772A4.902 4.902 0 015.45 2.525c.636-.247 1.363-.416 2.427-.465C8.901 2.013 9.256 2 11.685 2h.63zm-.081 1.802h-.468c-2.456 0-2.784.011-3.807.058-.975.045-1.504.207-1.857.344-.467.182-.8.398-1.15.748-.35.35-.566.683-.748 1.15-.137.353-.3.882-.344 1.857-.047 1.023-.058 1.351-.058 3.807v.468c0 2.456.011 2.784.058 3.807.045.975.207 1.504.344 1.857.182.466.399.8.748 1.15.35.35.683.566 1.15.748.353.137.882.3 1.857.344 1.054.048 1.37.058 4.041.058h.08c2.597 0 2.917-.01 3.96-.058.976-.045 1.505-.207 1.858-.344.466-.182.8-.398 1.15-.748.35-.35.566-.683.748-1.15.137-.353.3-.882.344-1.857.048-1.055.058-1.37.058-4.041v-.08c0-2.597-.01-2.917-.058-3.96-.045-.976-.207-1.505-.344-1.858a3.097 3.097 0 00-.748-1.15 3.098 3.098 0 00-1.15-.748c-.353-.137-.882-.3-1.857-.344-1.023-.047-1.351-.058-3.807-.058zM12 6.865a5.135 5.135 0 110 10.27 5.135 5.135 0 010-10.27zm0 1.802a3.333 3.333 0 100 6.666 3.333 3.333 0 000-6.666zm5.338-3.205a1.2 1.2 0 110 2.4 1.2 1.2 0 010-2.4z"
                  clip-rule="evenodd"
                ></path>
              </svg>
            </div>

            <!-- Hidden Space for Small Screens -->
            <div class="md:hidden w-18"></div>
          </div>
        </div>
      </div>

      <div
        class="px-5 py-2 flex gap-5 justify-center w-max mx-auto items-center border rounded-full border-gray-300"
      >
        <div class="flex items-center gap-2">
          <div
            class="w-0 h-0 -left-4 relative border-l-[5px] border-l-transparent border-b-[10px] border-b-orange-500 border-r-[5px] border-r-transparent"
          ></div>
          <div class="text-sm text-gray-700">SP</div>
        </div>
        <div class="bg-gray-100 h-4 w-1"></div>
        <div class="flex items-center gap-2">
          <div
            class="w-2 h-2 bg-indigo-600 rounded-full relative -left-4"
          ></div>
          <div class="text-sm text-gray-700">RI</div>
        </div>
      </div>
      <div>
        <div
          class="relative pt-6 mt-8 pb-4 text-center border-r border-l border-l-gray-300 border-r-gray-300"
        >
          <div class="mb-4 font-bold text-lg text-gray-700 relative px-8">
            <div
              style="top:4px;left:-2px;height:16px"
              class="absolute rounded-full grid-pill w-1 bg-white z-10"
            ></div>
            <div
              style="top:4px;right:-2px;height:16px"
              class="absolute rounded-full grid-pill w-1 bg-white z-10"
            ></div>
            Commitment Recommendations &amp; Management
          </div>
          <div class="text-sm text-gray-400 max-w-460 pb-4 mx-auto px-4">
            Maximize your savings with Reserved Instances and Savings Plans
            while mitigating long-term commitment risk.
          </div>
        </div>
      </div>
    </div>
  </div>

  <Seperator />
</section>

<section class="relative bg-black">
  <Pricing />
</section>

<style>
  .canvas-small-screen {
    display: flex; /* Default to hide on large screens */
  }

  /* Media query for small screens */
  @media (max-width: 640px) {
    .canvas-small-screen {
      display: block; /* Show canvas on small screens */
      height: 398.4px;
      width: 398.4px;
    }
  }
  .savings-optimization .dotted-background {
    -webkit-mask: linear-gradient(to bottom, transparent 0%, #fff 20%);
    mask: linear-gradient(to bottom, transparent 0%, #fff 20%);
  }

  .savings-optimization .savings-optimization-antimetal-border {
    -webkit-mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    mask-composite: exclude;
    padding: 1px;
  }
  /* Add this to your styles.css or equivalent */

  .bg-neutral-opaque-8 {
    background-color: rgba(0, 0, 0, 0.08); /* Example for opacity */
  }

  .text-body-small-regular {
    font-size: 0.875rem; /* Adjust as needed */
    font-weight: 400;
  }

  .text-body-small-medium {
    font-size: 0.875rem; /* Adjust as needed */
    font-weight: 500;
  }

  .dotted-background {
    background: radial-gradient(
      circle,
      rgba(0, 0, 0, 0.05) 20%,
      rgba(0, 0, 0, 0) 30%
    );
    background-size: 8px 8px; /* Smaller dots */
    background-position:
      0 0,
      4px 4px; /* Adjust position for diagonal effect */
    background-repeat: repeat;
    -webkit-mask: radial-gradient(circle, transparent 20%, #fff 30%);
    mask: radial-gradient(circle, transparent 20%, #fff 30%);
    -webkit-mask-composite: source-in;
    mask-composite: intersect;
  }

  .logos-container {
    display: grid;
    grid-template-columns: repeat(3, 2fr);
    grid-template-rows: repeat(4, auto);
    gap: 2em; /* Adjust the gap between items if necessary */
  }

  .logo-item {
    width: 5rem;
    height: 5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    border-radius: 50%;
  }

  .logo-label {
    position: absolute;
    bottom: -20px; /* Adjust as needed */
    left: 50%;
    transform: translateX(-50%);
    text-align: center;
    font-size: 10px;
    color: #999;
  }

  .savings-optimization .savings-optimization-antimetal-border::before {
    content: "";
    background: conic-gradient(
        from 5deg at 50% 50%,
        transparent 260deg,
        rgb(73, 201, 176) 360deg,
        transparent 361deg
      ),
      var(--neutral-opaque-8);
    height: 120px;
    width: 120px;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%) rotate(0deg);
    animation: savings-optimization-antimetal 4s linear infinite;
  }

  @keyframes savings-optimization-antimetal {
    to {
      transform: translate(-50%, -50%) rotate(360deg);
    }
  }

  .savings-optimization .savings-optimization-antimetal-inner {
    box-shadow:
      0 0 0 1px #0e3f7e0f,
      0 1px 1px -0.5px #2a334608,
      0 3px 3px -1.5px #2a33460d,
      0 5px 5px -2.5px #2a334608,
      0 16px 16px -8px #2a33460a;
  }

  .savings-optimization .savings-optimization-antimetal-dot {
    box-shadow: 0 1px 2px #0e3f7e14 inset;
  }

  section {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    flex: 0.6;
    height: auto;
    position: relative;
  }
  @keyframes customAnimation {
    0% {
      transform: scale(0.5);
      opacity: 0;
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }

  @keyframes customAnimation {
    0% {
      transform: scale(0.5);
      opacity: 0;
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }

  .custom-animation {
    animation: customAnimation 0.6s ease-in-out;
  }
  .grid-item {
    width: 10.5rem;
    height: 10.38rem;
    /* padding: 5rem 4rem 5rem 4rem; */
    border-radius: 1rem; /* rounded-2xl */
    /* border: 1px solid #bbcbff; Adjust color if necessary */
    box-shadow: 0 4px 6px rgba(98, 98, 98, 0.1);
    grid-column: span 1 / span 1;
  }

  /* Scoped styles for this component */

  /* Transferring your custom styles */
  .problem-items {
    -webkit-mask-image: linear-gradient(
      to right,
      transparent 0%,
      #d9d9d9 20%,
      #d9d9d9 80%,
      transparent 100%
    );
    mask-image: linear-gradient(
      to right,
      transparent 0%,
      #d9d9d9 20%,
      #d9d9d9 80%,
      transparent 100%
    );
  }

  @media (max-width: 767px) {
    .problem-items {
      -webkit-mask-image: linear-gradient(
        to right,
        transparent 0%,
        #d9d9d9 10%,
        #d9d9d9 90%,
        transparent 100%
      );
      mask-image: linear-gradient(
        to right,
        transparent 0%,
        #d9d9d9 10%,
        #d9d9d9 90%,
        transparent 100%
      );
    }

    .problem-item:nth-child(n + 9) {
      display: none;
      animation: unset;
    }
  }

  .problem-item {
    animation: problem-item-box 7s infinite backwards;
    transform: translateZ(0);
    box-shadow:
      0 0 0 1px #0e3f7e0f,
      0 1px 1px -0.5px #2a334608,
      0 2px 2px -1px #2a33460a,
      0 3px 3px -1.5px #2a33460a,
      0 5px 5px -2.5px #2a334608,
      0 10px 10px -5px #2a334608,
      0 24px 24px -8px #2a334608;
    animation-delay: var(--animation-delay);
    background-color: var(--white-24);
  }

  @keyframes problem-item-box {
    17.8%,
    57.1% {
      background-color: var(--white-24);
    }

    25%,
    50% {
      background-color: var(--white-72);
    }
  }

  .problem-item-label {
    animation: problem-item-label 7s infinite backwards;
    opacity: 0;
    transform: translateY(10px);
    animation-delay: calc(var(--animation-delay) + 0.1s);
  }

  @keyframes problem-item-label {
    0%,
    17.8% {
      opacity: 0;
      transform: translateY(10px);
    }

    25%,
    78.5% {
      opacity: 1;
      transform: translateY(0);
      filter: blur(0);
    }

    85.6% {
      filter: blur(10px) grayscale(1);
      opacity: 0;
      transform: translateY(-20px);
    }
  }

  .problem-item-icon {
    animation: problem-item-icon 7s infinite backwards;
    opacity: 0;
    transform: translateY(20px);
    filter: blur(10px);
    animation-delay: var(--animation-delay);
    background-color: var(--background-color-p3, var(--background-color-hex));
  }

  @keyframes problem-item-icon {
    0%,
    14.2% {
      opacity: 0;
      transform: translateY(20px);
      filter: blur(10px) grayscale(1);
    }

    21.4%,
    78.5% {
      opacity: 1;
      transform: translateY(0);
      filter: blur(0);
    }

    82% {
      filter: blur(5px) grayscale(1);
    }

    to {
      opacity: 0;
      transform: translateY(-100px) scale(2);
      filter: blur(20px) grayscale(1);
    }
  }

  .problem-item-icon-img {
    animation: problem-item-icon-img 7s infinite backwards;
    animation-delay: var(--animation-delay);
    height: 0px;
    width: 0px;
    overflow: hidden;
  }

  @keyframes problem-item-icon-img {
    0%,
    17.8%,
    83% {
      height: 0px;
      width: 0px;
      transform: translateY(0);
    }

    32.1%,
    78.5% {
      height: 32px;
      width: 32px;
      opacity: 1;
      transform: translateY(0);
    }

    82% {
      opacity: 0;
      height: 32px;
      width: 32px;
      transform: translateY(-10px);
    }
  }
  /* General Styles */
  .savings-optimization-logos-item {
    position: absolute;
    border-radius: inherit;
  }

  .savings-optimization-logos-item-icon:before,
  .savings-optimization-logos-item-icon:after {
    content: "";
    transition-property: all;
    transition-timing-function: cubic-bezier(0.6, 0.6, 0, 1);
    transition-duration: 0.45s;
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    border-radius: inherit;
  }

  /* Icon Shadows */
  .savings-optimization-logos-item-icon:before {
    box-shadow:
      0 0 0 1px #0e3f7e0f,
      0 1px 1px -0.5px #2a334608,
      0 3px 3px -1.5px #2a33460d,
      0 5px 5px -2.5px #2a334608,
      0 16px 16px -8px #2a33460a;
  }

  .savings-optimization-logos-item-icon:after {
    box-shadow:
      0 0 0 1px #00685b0f,
      0 1px 1px -0.5px #00685b08,
      0 3px 3px -1.5px #00ccaf1f,
      0 5px 5px -2.5px #00ccaf1f,
      0 16px 16px -8px #00ccaf29,
      0 24px 24px -8px #00ccaf1f;
    opacity: 0;
  }

  /* Item Positions */
  .savings-optimization-logos-item:nth-child(1) {
    left: 11.6666666667%;
    top: 15.9090909091%;
  }

  .savings-optimization-logos-item:nth-child(2) {
    left: 15.8333333333%;
    top: 51.7045454545%;
  }

  .savings-optimization-logos-item:nth-child(3) {
    top: 65.9090909091%;
    left: 50%;
    transform: translate(-50%);
  }

  .savings-optimization-logos-item:nth-child(4) {
    right: 15.8333333333%;
    top: 51.7045454545%;
  }

  .savings-optimization-logos-item:nth-child(5) {
    right: 11.6666666667%;
    top: 15.9090909091%;
  }

  /* Border and Animation */
  .savings-optimization-logos-item-border {
    -webkit-mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    mask-composite: xor;
    -webkit-mask-composite: xor;
    padding: 1px;
  }

  .savings-optimization-logos-item-border:before {
    content: "";
    background: conic-gradient(
      from 5deg at 50% 50%,
      transparent 260deg,
      rgb(73, 201, 176) 360deg,
      transparent 361deg
    );
    animation: savings-optimization-logos-border 1s linear infinite;
    height: 100px;
    width: 100px;
    transition-property: all;
    transition-timing-function: cubic-bezier(0.6, 0.6, 0, 1);
    transition-duration: 0.45s;
    left: 50%;
    transform: translate(-50%, -50%) rotate(var(--tw-rotate))
      skew(var(--tw-skew-x)) skewY(var(--tw-skew-y)) scaleX(var(--tw-scale-x))
      scaleY(var(--tw-scale-y));
  }

  @keyframes savings-optimization-logos-border {
    to {
      transform: translate(-50%, -50%) rotate(360deg);
    }
  }

  /* Optimized Label and Check Styles */
  .savings-optimization-logos-item-optimized-label {
    color: #008873b8;
    color: color(display-p3 0.1647 0.5255 0.4549 / 0.72);
  }

  .savings-optimization-logos-item-check,
  .savings-optimization-logos-item-optimized-check {
    background: linear-gradient(
        180deg,
        rgba(255, 255, 255, 0.2) 0%,
        rgba(255, 255, 255, 0.2) 100%
      ),
      #7fd8c6;
    background: linear-gradient(
        180deg,
        color(display-p3 1 1 1 / 0.2) 0%,
        color(display-p3 1 1 1 / 0) 100%
      ),
      color(display-p3 0.498 0.8471 0.7765);
    box-shadow:
      0 0 0 1px #00a4851f inset,
      0 0.5px 0.5px #00b89614,
      0 1.5px 1.5px #00b89614,
      0 2.5px 2.5px #00b89614,
      0 4px 4px #00b89614,
      0 8px 8px #00b8961a;
  }

  .savings-optimization-logos-item-check svg,
  .savings-optimization-logos-item-optimized-check svg {
    filter: drop-shadow(0px 0.5px 0.5px rgba(0, 184, 150, 0.24))
      drop-shadow(0px 1.5px 1.5px rgba(0, 184, 150, 0.24))
      drop-shadow(0px 2.5px 2.5px rgba(0, 184, 150, 0.24))
      drop-shadow(0px 4px 4px rgba(0, 184, 150, 0.24));
  }

  /* Step-Based Styles */
  .savings-optimization-logos-item[data-step="0"]
    .savings-optimization-logos-item-icon:before,
  .savings-optimization-logos-item[data-step="1"]
    .savings-optimization-logos-item-icon:before,
  .savings-optimization-logos-item[data-step="2"]
    .savings-optimization-logos-item-icon:before {
    opacity: 0;
  }

  .savings-optimization-logos-item[data-step="0"]
    .savings-optimization-logos-item-icon:after,
  .savings-optimization-logos-item[data-step="1"]
    .savings-optimization-logos-item-icon:after,
  .savings-optimization-logos-item[data-step="2"]
    .savings-optimization-logos-item-icon:after {
    opacity: 1;
  }

  .savings-optimization-logos-item[data-step="0"]
    .savings-optimization-logos-item-border-wrapper,
  .savings-optimization-logos-item[data-step="1"]
    .savings-optimization-logos-item-border-wrapper,
  .savings-optimization-logos-item[data-step="2"]
    .savings-optimization-logos-item-border-wrapper {
    --tw-scale-x: 1;
    --tw-scale-y: 1;
    transform: translate(var(--tw-translate-x), var(--tw-translate-y))
      rotate(var(--tw-rotate)) skew(var(--tw-skew-x)) skewY(var(--tw-skew-y))
      scaleX(var(--tw-scale-x)) scaleY(var(--tw-scale-y));
    opacity: 1;
    --tw-grayscale: grayscale(0);
    filter: var(--tw-blur) var(--tw-brightness) var(--tw-contrast)
      var(--tw-grayscale) var(--tw-hue-rotate) var(--tw-invert)
      var(--tw-saturate) var(--tw-sepia) var(--tw-drop-shadow);
  }

  .savings-optimization-logos-item[data-step="1"]
    .savings-optimization-logos-item-border-wrapper
    .savings-optimization-logos-item-border:before,
  .savings-optimization-logos-item[data-step="2"]
    .savings-optimization-logos-item-border-wrapper
    .savings-optimization-logos-item-border:before {
    opacity: 0;
  }

  .savings-optimization-logos-item[data-step="1"]
    .savings-optimization-logos-item-optimized,
  .savings-optimization-logos-item[data-step="2"]
    .savings-optimization-logos-item-optimized,
  .savings-optimization-logos-item[data-step="3"]
    .savings-optimization-logos-item-optimized {
    --tw-translate-y: 0px;
    --tw-scale-x: 1;
    --tw-scale-y: 1;
    transform: translate(var(--tw-translate-x), var(--tw-translate-y))
      rotate(var(--tw-rotate)) skew(var(--tw-skew-x)) skewY(var(--tw-skew-y))
      scaleX(var(--tw-scale-x)) scaleY(var(--tw-scale-y));
    opacity: 1;
    --tw-grayscale: grayscale(0);
    filter: var(--tw-blur) var(--tw-brightness) var(--tw-contrast)
      var(--tw-grayscale) var(--tw-hue-rotate) var(--tw-invert)
      var(--tw-saturate) var(--tw-sepia) var(--tw-drop-shadow);
  }

  .savings-optimization-logos-item[data-step="2"]
    .savings-optimization-logos-item-optimized {
    box-shadow:
      0 0 0 1px #0e3f7e0f,
      0 1px 1px -0.5px #0e3f7e08,
      0 3px 3px -1.5px #0e3f7e0d,
      0 5px 5px -2.5px #0e3f7e08,
      0 16px 16px -8px #0e3f7e0a,
      0 24px 24px -8px #0e3f7e0f;
    height: 28px;
    width: 106px;
    --tw-translate-y: -8px;
    background-color: var(--white-72);
    --tw-backdrop-blur: blur(10px);
    -webkit-backdrop-filter: var(--tw-backdrop-blur)
      var(--tw-backdrop-brightness) var(--tw-backdrop-contrast)
      var(--tw-backdrop-grayscale) var(--tw-backdrop-hue-rotate)
      var(--tw-backdrop-invert) var(--tw-backdrop-opacity)
      var(--tw-backdrop-saturate) var(--tw-backdrop-sepia);
    backdrop-filter: var(--tw-backdrop-blur) var(--tw-backdrop-brightness)
      var(--tw-backdrop-contrast) var(--tw-backdrop-grayscale)
      var(--tw-backdrop-hue-rotate) var(--tw-backdrop-invert)
      var(--tw-backdrop-opacity) var(--tw-backdrop-saturate)
      var(--tw-backdrop-sepia);
  }

  .savings-optimization-logos-item[data-step="2"]
    .savings-optimization-logos-item-optimized-label,
  .savings-optimization-logos-item[data-step="3"]
    .savings-optimization-logos-item-optimized-label {
    --tw-translate-x: 0px;
    transform: translate(var(--tw-translate-x), var(--tw-translate-y))
      rotate(var(--tw-rotate)) skew(var(--tw-skew-x)) skewY(var(--tw-skew-y))
      scaleX(var(--tw-scale-x)) scaleY(var(--tw-scale-y));
    opacity: 1;
    transition-delay: 0.2s;
  }

  .savings-optimization-logos-item[data-step="2"]
    .savings-optimization-logos-item-optimized-check,
  .savings-optimization-logos-item[data-step="3"]
    .savings-optimization-logos-item-optimized-check {
    --tw-translate-x: 0px;
    transform: translate(var(--tw-translate-x), var(--tw-translate-y))
      rotate(var(--tw-rotate)) skew(var(--tw-skew-x)) skewY(var(--tw-skew-y))
      scaleX(var(--tw-scale-x)) scaleY(var(--tw-scale-y));
  }

  .savings-optimization-logos-item[data-step="3"]
    .savings-optimization-logos-item-optimized {
    --tw-translate-y: -32px;
    opacity: 0;
    --tw-blur: blur(6px);
    --tw-grayscale: grayscale(100%);
    filter: var(--tw-blur) var(--tw-brightness) var(--tw-contrast)
      var(--tw-grayscale) var(--tw-hue-rotate) var(--tw-invert)
      var(--tw-saturate) var(--tw-sepia) var(--tw-drop-shadow);
  }

  .savings-optimization-logos-item[data-step="3"]
    .savings-optimization-logos-item-check {
    --tw-translate-y: 0px;
    opacity: 1;
    --tw-blur: blur(0px);
    --tw-grayscale: grayscale(0);
    filter: var(--tw-blur) var(--tw-brightness) var(--tw-contrast)
      var(--tw-grayscale) var(--tw-hue-rotate) var(--tw-invert)
      var(--tw-saturate) var(--tw-sepia) var(--tw-drop-shadow);
  }
</style>
