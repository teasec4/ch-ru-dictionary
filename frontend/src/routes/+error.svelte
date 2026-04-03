<script lang="ts">
    import { page } from "$app/state";
    import { Home, ArrowLeft } from "lucide-svelte";

    let { data } = $props();
    const translations = data.translations;
    
    let error = $derived(page.error);
</script>

<div class="min-h-screen w-full flex flex-col items-center justify-center p-4">
    <div class="text-center max-w-md">
        <h1 class="text-6xl font-bold text-accent mb-4">404</h1>
        <h2 class="text-xl font-semibold text-text-primary mb-2">
            {translations.pageNotFound || "Страница не найдена"}
        </h2>
        <p class="text-dict-2 mb-6">
            {translations.pageNotFoundDesc || "Извините, запрашиваемая страница не существует."}
        </p>
        
        <div class="flex gap-3 justify-center">
            <a 
                href="/{page.params.lang || 'ru'}" 
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-dict-1 transition"
            >
                <Home class="w-4 h-4" />
                {translations.home || "На главную"}
            </a>
            
            <button 
                onclick={() => history.back()}
                class="flex items-center gap-2 px-4 py-2 border border-dict-3 text-dict-1 rounded-lg hover:bg-dict-4 transition"
            >
                <ArrowLeft class="w-4 h-4" />
                {translations.back || "Назад"}
            </button>
        </div>
    </div>
</div>