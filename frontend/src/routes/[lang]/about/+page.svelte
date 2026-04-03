<script lang="ts">
    import { page } from "$app/state";
    let { data } = $props();
    const t = data.translations;
</script>

<svelte:head>
    <title>{t.aboutTitle} - {t.siteTitle}</title>
    <meta name="description" content={t.aboutText} />
</svelte:head>

<div class="w-full max-w-2xl mx-auto py-8 px-4">
    <h1 class="text-2xl sm:text-3xl font-bold text-text-primary mb-6">{t.aboutTitle}</h1>
    
    <div class="space-y-6 text-dict-1">
        <p class="text-base sm:text-lg leading-relaxed">{t.aboutText}</p>
        
        <div class="p-4 rounded-xl bg-dict-4/20 border border-dict-4">
            <h2 class="text-lg font-semibold text-text-primary mb-2">{t.missionTitle}</h2>
            <p class="text-dict-2">{t.missionText}</p>
        </div>
    </div>
</div>