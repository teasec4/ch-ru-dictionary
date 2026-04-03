<script lang="ts">
    let { data } = $props();
    const t = data.translations;
</script>

<svelte:head>
    <title>{t.privacyTitle} - {t.siteTitle}</title>
</svelte:head>

<div class="w-full max-w-2xl mx-auto py-8 px-4">
    <h1 class="text-2xl sm:text-3xl font-bold text-text-primary mb-6">{t.privacyTitle}</h1>
    
    <div class="prose prose-dict text-dict-2 space-y-4">
        <p>{t.privacyText}</p>
    </div>
</div>