<script lang="ts">
  import { settingsStructure } from '$lib/config/settings.config';
  import { page } from '$app/stores';
</script>

<nav class="h-full pr-4 py-6 border-r border-border overflow-y-auto">
  {#each settingsStructure as section}
    <div class="mb-8">
      <h3 class="text-base font-semibold mb-2 text-foreground">{section.label}</h3>
      <ul class="space-y-1">
        {#each section.children as tab}
          <li>
            <a
              href={tab.path.replace('[shop]', $page.params.shop)}
              class="nav-item"
              class:nav-item-active={ $page.url.pathname.startsWith(tab.path.replace('[shop]', $page.params.shop)) }
              class:nav-item-inactive={!$page.url.pathname.startsWith(tab.path.replace('[shop]', $page.params.shop)) }
              aria-current={ $page.url.pathname.startsWith(tab.path.replace('[shop]', $page.params.shop)) ? 'page' : undefined }
            >
              {tab.label}
            </a>
          </li>
        {/each}
      </ul>
    </div>
  {/each}
</nav> 