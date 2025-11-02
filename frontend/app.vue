<script setup lang="ts">
import { useAsyncData } from '#app'
import Button from './components/ui/Button.vue'
import { cn } from './lib/utils'
import { LoaderCircle, Moon, Sun } from 'lucide-vue-next'
import { onBeforeMount, ref } from 'vue'

const params = new URLSearchParams(window.location.search)

enum GraphDuration {
	Day = 'day',
	Week = 'week',
	Month = 'month',
	Year = 'year',
	All = 'all'
}

const durationLabel = {
	[GraphDuration.Day]: '24h',
	[GraphDuration.Week]: '7d',
	[GraphDuration.Month]: '30d',
	[GraphDuration.Year]: '1y',
	[GraphDuration.All]: 'all'
}

interface UptimeResponse {
	Overall: 'all_good' | 'degraded' | 'down' | 'partial' | 'unknown'
	Details: Array<{
		iconUrl: string
		title: string
		slug: string
		url?: string
		status: 'up' | 'down' | 'degrade' | 'unknown'

		responseOverall: string
		response24h: string
		response7d: string
		response30d: string
		response1y: string

		uptimeOverall: string
		uptime24h: string
		uptime7d: string
		uptime30d: string
		uptime1y: string
	}>
}

const overallLabel = {
	all_good: '✅ All systems operational',
	degraded: '🐌 Degraded performance',
	down: '❌ Complete outage',
	partial: '⚠️ Partial outage',
	unknown: '❓ Unknown'
}

const {
	data: apiData,
	status: apiStatus,
	error: apiErr
} = useAsyncData<UptimeResponse>(async () => {
	const user = params.get('user')
	if (!user || user === '') throw new Error('Username is required')

	const repo = params.get('repo')
	if (!repo || repo === '') throw new Error('Repository is required')

	const resp = await fetch(`/api/${user}/${repo}`)

	if (!resp.ok) throw new Error("Can't fetch data")

	// eslint-disable-next-line @typescript-eslint/no-unsafe-return
	return resp.json()
})

const preferredColorScheme = window.matchMedia('(prefers-color-scheme: dark)')
	.matches
	? 'dark'
	: window.matchMedia('(prefers-color-scheme: light)').matches
		? 'light'
		: null

const configuredColorScheme =
	localStorage.getItem('color-scheme') === 'dark'
		? 'dark'
		: localStorage.getItem('color-scheme') === 'light'
			? 'light'
			: null

const colorScheme = ref<'dark' | 'light'>(
	configuredColorScheme ?? preferredColorScheme ?? 'dark'
)

function toggleColorScheme(): void {
	localStorage.setItem('color-scheme', colorScheme.value)
	if (colorScheme.value === 'light') {
		document.querySelector('html')?.classList.add('dark')
		localStorage.setItem('color-scheme', 'dark')
		colorScheme.value = 'dark'
	} else {
		document.querySelector('html')?.classList.remove('dark')
		localStorage.setItem('color-scheme', 'light')
		colorScheme.value = 'light'
	}
}

onBeforeMount(() => {
	if (colorScheme.value === 'dark') {
		document.querySelector('html')?.classList.add('dark')
	}
})

const viewDuration = ref<GraphDuration>(GraphDuration.All)
</script>

<template>
	<div class="mx-auto w-11/12 max-w-3xl">
		<div v-if="apiErr" class="grid h-dvh place-items-center">
			{{ apiErr.message }}
		</div>
		<div
			v-if="apiStatus === 'pending'"
			class="grid h-dvh place-items-center"
		>
			<LoaderCircle class="animate-spin" :size="30" />
		</div>
		<div v-if="apiData" class="my-7 flex flex-col gap-3">
			<div class="flex justify-between text-3xl font-semibold">
				{{ overallLabel[apiData.Overall] }}
				<Button
					variant="outline-solid"
					size="icon"
					@click="toggleColorScheme"
				>
					<Sun v-if="colorScheme === 'light'" class="size-4" />
					<Moon v-else class="size-4" />
				</Button>
			</div>

			<div class="ml-auto grid w-full max-w-72 grid-cols-5 gap-2">
				<Button
					v-for="duration in Object.values(GraphDuration)"
					:key="duration"
					:disabled="viewDuration === duration"
					size="xs"
					variant="secondary"
					class="bg-card shadow-sm transition-all"
					@click="viewDuration = duration"
				>
					<span class="mb-px text-xs">
						{{ durationLabel[duration] }}
					</span>
				</Button>
			</div>

			<div
				v-for="detail in apiData.Details"
				:key="detail.slug"
				class="grid h-20 grid-cols-[.4rem_1fr] gap-4 overflow-hidden rounded-md border"
				:style="{
					backgroundImage: `url(/api/graph/${params.get('user')}/${params.get('repo')}/${detail.slug}/${viewDuration.toString()})`,
					backgroundSize: 'contain',
					backgroundPosition: 'center right',
					backgroundRepeat: 'no-repeat'
				}"
			>
				<div
					:class="
						cn(
							colorScheme === 'dark' && {
								'bg-green-600': detail.status === 'up',
								'bg-red-700': detail.status === 'down',
								'bg-yellow-500': detail.status === 'degrade',
								'bg-gray-500': detail.status === 'unknown'
							},
							colorScheme === 'light' && {
								'bg-green-400': detail.status === 'up',
								'bg-red-500': detail.status === 'down',
								'bg-yellow-300': detail.status === 'degrade',
								'bg-gray-300': detail.status === 'unknown'
							}
						)
					"
				/>

				<div class="flex flex-col justify-center gap-2">
					<div
						:class="
							cn(
								'flex flex-row items-center gap-2 font-medium',
								colorScheme === 'dark' && {
									'text-green-400': detail.status === 'up',
									'animate-bounce text-red-500':
										detail.status === 'down',
									'text-yellow-400':
										detail.status === 'degrade',
									'text-gray-400': detail.status === 'unknown'
								},
								colorScheme === 'light' && {
									'text-green-600': detail.status === 'up',
									'animate-bounce text-red-600':
										detail.status === 'down',
									'text-yellow-600':
										detail.status === 'degrade',
									'text-gray-600': detail.status === 'unknown'
								}
							)
						"
					>
						<img :src="detail.iconUrl" class="size-5" />
						<span class="mb-1">{{ detail.title }}</span>
					</div>
					<div
						class="flex gap-2 text-xs text-secondary-foreground/70"
					>
						<code
							>🟢
							{{
								(() => {
									switch (viewDuration) {
										case GraphDuration.Day:
											return detail.uptime24h
										case GraphDuration.Week:
											return detail.uptime7d
										case GraphDuration.Month:
											return detail.uptime30d
										case GraphDuration.Year:
											return detail.uptime1y
										case GraphDuration.All:
											return detail.uptimeOverall
									}
								})()
							}}</code
						>
						<code
							>⏳
							{{
								(() => {
									switch (viewDuration) {
										case GraphDuration.Day:
											return detail.response24h
										case GraphDuration.Week:
											return detail.response7d
										case GraphDuration.Month:
											return detail.response30d
										case GraphDuration.Year:
											return detail.response1y
										case GraphDuration.All:
											return detail.responseOverall
									}
								})()
							}}ms</code
						>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
