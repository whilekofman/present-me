<template>
    <div class="mx-auto">
        <div class="bg-gradient-to-b from-gray-800 to-black text-white font-mono text-sm text-center py-2 shadow">
          {{ $route.params.org }}/{{ $route.params.repo }}/pull/{{ $route.params.pull }}/review-{{ $route.params.review }}
        </div>
        <div v-if="pending" class="py-10">
          <div class="animate-pulse mx-auto text-center text-4xl font-bold">
            <GradientText>Loading...</GradientText>
          </div>
        </div>
        <div v-else-if="error" class="mx-auto max-w-3xl py-10">
          <div class="bg-orange-100 px-2 pb-2">
            <div class="text-xs text-center underline py-4">{{ error }}</div>
            <div class="bg-white p-4 text-center border border-orange-200 rounded">
              <code>{{ error.data }}</code>
            </div>
          </div>
        </div>
        <div v-else>
          <Review :model="data" />
        </div>
    </div>
</template>

<script setup lang="ts">
useHead({
  title: 'present-me',
});
const route = useRoute();
const { pending, data, error } = await useFetch('/api/review', {
  params: route.params,
  server: false,
  initialCache: false,
  transform: v => JSON.parse(v),
});
</script>
