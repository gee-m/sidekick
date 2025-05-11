# Life Tracking System

The goal is to make tracking "stupidly easy" - as easy as pressing a button or sending a quick message, or maybe even snapping a pic. Some questions are things that should be asked daily (e.g., how you feel), some should be inputted by the user but may need a reminder if not inputted during the day.

## Project Overview

A comprehensive life tracking system designed to collect, correlate, and analyze various aspects of daily life to identify patterns affecting happiness and wellbeing. Inspired by Bryan Johnson's tracking work and Andrej Karpathy's sleep tracking analysis.

### Core Philosophy

> "My sleep scores correlate strongly with the quality of work I am able to do that day. When my score is low, I lack agency, I lack courage, I lack creativity, I'm simply tired. When my sleep score is high, I can power through anything." - Andrej Karpathy

This system aims to identify similar correlations across multiple life domains to optimize wellbeing through data-driven insights.

## Core Objectives

1. **Low Friction Data Collection** - Make tracking "stupidly easy" through mobile-optimized interfaces
2. **Comprehensive Life Domains** - Track sleep, nutrition, exercise, work, relationships, and more
3. **Meaningful Correlations** - Identify connections between behaviors and wellbeing outcomes
4. **Actionable Insights** - Generate practical recommendations based on personal data patterns

## Domain Model

> **Event-Based Architecture Note**: While the model below describes aggregate data structures, the underlying implementation should follow an event-based architecture. Events (like "TookSupplement", "WashedHair", "AteGluten") would be the primary data entities, with these domain models generated as aggregated views. This approach enables more flexible analysis, better historical tracking, and cleaner data management.

### DailyRecord
- `date`: Date - The date of this record
- `subjective`: SubjectiveMetrics - Subjective wellbeing indicators
- `sleep`: SleepData - Sleep metrics including Oura data
- `nutrition`: NutritionData - Food, drink, and meal timing data
- `activities`: ActivityData - Exercise, nature time, and other activities
- `relationships`: RelationshipData - Social and relationship interactions
- `supplements`: SupplementData - Supplement and medication tracking
- `timeUsage`: TimeUsageData - How time was allocated across activities
- `selfCare`: SelfCareData - Self-care activities and tracking
- `work`: WorkData - Work-related metrics and productivity
- `derivedMetrics`: DerivedMetrics - Automatically calculated metrics (streaks, days since, etc.)
- `environmentalData`: EnvironmentalData - External environmental factors

### QuickAction
- `id`: String - Unique identifier
- `type`: Enum(SUPPLEMENT, MEAL, ACTIVITY, SELFCARE, etc.) - Type of action
- `label`: String - Display text for the button
- `payload`: JSON - Data to be submitted when tapped
- `icon`: String - Icon identifier
- `usageCount`: Number - How many times action has been used
- `lastUsed`: DateTime - When action was last used
- `contextTags`: [String] - When this action is most relevant (morning, evening, etc.)

### DerivedMetrics
- `streaks`: Map<ActivityType, Number> - Consecutive days of activities/habits
- `daysSinceEvents`: Map<EventType, Number> - Days since various tracked events
  - Contains all "days since" metrics (work, free day, holiday, fast food, gluten,
    family contact, exercise, supplements, nature time, animal interaction, etc.)
  - Also tracks days since hygiene events (washed hair, washed body, etc.)
  - Also tracks days since housekeeping events (bed sheets cleaned, laundry done, etc.)
  - **Note**: These metrics would be dynamically calculated from event logs rather than stored directly
- `consistencyScores`: Map<MetricType, Score> - Consistency of various metrics
- `correlationInsights`: [CorrelationInsight] - Automatically detected correlations

### SubjectiveMetrics
- `morningRestedness`: Score(1-10) - How rested upon waking
- `morningHappiness`: Score(1-10) - Happiness level in morning
- `morningAnxiety`: Score(1-10) - Anxiety level in morning
- `eveningHappiness`: Score(1-10) - Happiness level in evening
- `eveningAnxiety`: Score(1-10) - Anxiety level in evening
- `productivityFeeling`: Score(1-10) - Subjective productivity rating
- `lifeAdminPressure`: Score(1-10) - Pressure from life admin tasks

### SleepData
- `ouraData`: OuraMetrics - Data imported from Oura API
  - `sleepScore`: Number - Overall sleep score
  - `deepSleep`: Duration - Time in deep sleep
  - `remSleep`: Duration - Time in REM sleep
  - `lightSleep`: Duration - Time in light sleep
  - `awake`: Duration - Time awake during sleep period
  - `totalSleep`: Duration - Total sleep duration
  - `restingHeartRate`: Number - RHR during sleep
  - `heartRateVariability`: Number - HRV during sleep
  - `respiratoryRate`: Number - Breathing rate during sleep
  - `bodyTemperature`: Number - Body temperature deviation
- `bedTime`: Time - Time went to bed
- `wakeTime`: Time - Time woke up
- `bedTimeConsistency7d`: Duration - 7-day consistency of bedtime
- `bedTimeConsistency30d`: Duration - 30-day consistency of bedtime
- `wakeTimeConsistency7d`: Duration - 7-day consistency of wake time
- `wakeTimeConsistency30d`: Duration - 30-day consistency of wake time
- `hoursBeforeSleepLastMeal`: Number - Hours between last meal and sleep
- `sleepConditions`: SleepConditions - Environment and sleep aids used

### SleepConditions
- `sleepMask`: Boolean - Whether sleep mask was used
- `earplugs`: Boolean - Whether earplugs were used
- `sleepToPodcast`: Boolean - Whether fell asleep to podcast/audio
- `sleepToYoutube`: Boolean - Whether fell asleep to YouTube videos
- `greyNoise`: Boolean - Whether grey/white noise machine was used
- `noisePollution`: Boolean - Whether noise disturbance was present
- `readBeforeBed`: Boolean - Whether reading occurred before sleep
- `noseStrips`: Boolean - Whether nose strips were used
- `mouthTape`: Boolean - Whether mouth tape was used
- `withPartner`: Boolean - Whether slept with partner
- `hotShowerBefore`: Boolean - Whether took hot shower before bed
- `showeredBeforeBed`: Boolean - Whether showered before bed
- `bathedBeforeBed`: Boolean - Whether bathed before bed
- `lightPollution`: Boolean - Whether light disturbance was present
- `hungryBeforeBed`: Boolean - Whether felt hungry at bedtime
- `thirstyDuringNight`: Boolean - Whether woke up thirsty during the night
- `airConditioned`: Boolean - Whether air conditioning was used
- `groundingUsed`: Boolean - Whether grounding mat/sheet was used
- `alcoholConsumed`: Boolean - Whether alcohol was consumed before bed
- `screenExposureBefore`: Boolean - Whether screens were used shortly before bed
- `selfGratificationBefore`: Boolean - Whether self-gratification occurred before sleep
- `selfGratificationWithContent`: Boolean - Whether self-gratification involved content consumption
- `woreSocks`: Boolean - Whether socks were worn during sleep
- `woreShirt`: Boolean - Whether a shirt was worn during sleep
- `roomTemperature`: String/Number - Temperature of sleep environment
- `mattressType`: String - Type of mattress used
- `pillowType`: String - Type of pillow used

### NutritionData
- `meals`: [Meal] - List of meals consumed
  - `mealType`: Enum(BREAKFAST, LUNCH, DINNER, SNACK)
  - `time`: Time - When meal was consumed
  - `category`: [FoodCategory] - Types of food consumed
  - `satiety`: Score(1-10) - How filling the meal was (10 = overate)
  - `glutenConsumed`: Boolean - Whether meal contained gluten
  - `ingredients`: [String] - Key ingredients in the meal
- `fastFoodConsumed`: Boolean - Whether fast food was consumed
- `alcoholConsumed`: Boolean - Whether alcohol was consumed
- `caffeineConsumed`: Boolean - Whether caffeine was consumed
- `waterIntake`: Number - Estimated water consumption in ml

### ActivityData
- `stepCount`: Number - Number of steps taken
- `exercisePerformed`: Boolean - Whether exercise was performed
- `exerciseType`: [ExerciseType] - Types of exercise performed
- `exerciseDuration`: Duration - Duration of exercise
- `outdoorTime`: Duration - Time spent outdoors
- `natureTime`: Duration - Time spent in nature

### RelationshipData
- `partnerProximity`: Boolean - Whether partner was physically present
- `partnerInteraction`: Score(1-10) - Quality of partner interactions
- `intimacyWithPartner`: Boolean - Whether intimacy occurred with partner
- `conflictWithPartner`: Boolean - Whether conflict occurred with partner
- `conflictIntensity`: Score(1-10) - Intensity of any conflict
- `familyContact`: Boolean - Whether family contact occurred
- `socialInteraction`: Boolean - Whether social interaction occurred
- `animalInteraction`: Boolean - Whether interaction with pets/animals occurred

### SupplementData
- `supplementsTaken`: [Supplement] - Supplements consumed
  - `name`: String - Name of supplement
  - `dose`: String - Dose of supplement
  - `time`: Time - Time supplement was taken
- `creatineAmount`: Number - Amount of creatine consumed (mg)
- `dutasterideTaken`: Boolean - Whether dutasteride was taken
- `minoxidilApplied`: Boolean - Whether minoxidil was applied
- `otherMedications`: [Medication] - Other medications taken

### TimeUsageData
- `workScreenTime`: Duration - Screen time for work
- `personalComputerTime`: Duration - Personal computer usage time
- `phoneScreenTime`: Duration - Phone screen time
- `phoneAppBreakdown`: Map<App, Duration> - Breakdown of app usage
  - Including Instagram, YouTube, etc.
- `socialMedia`: SocialMediaUsage - Detailed social media consumption
- `workPerformed`: Boolean - Whether work was performed
- `personalProjectWork`: Boolean - Whether personal projects were worked on
- `workProductivity`: Score(1-10) - Subjective work productivity

### SocialMediaUsage
- `instagramReels`: String/Duration - Time spent watching Instagram reels
- `tiktok`: String/Duration - Time spent on TikTok
- `youtube`: String/Duration - Time spent on YouTube
- `twitter`: String/Duration - Time spent on Twitter/X
- `facebook`: String/Duration - Time spent on Facebook
- `reddit`: String/Duration - Time spent on Reddit
- `other`: Map<Platform, Duration> - Other social platforms

### SelfCareData
- `selfCareActivities`: [SelfCareActivity] - Self-care activities performed
  - `type`: Enum(HAIRCUT, BATH, MASSAGE, SKINCARE, MEDITATION, OTHER)
  - `duration`: Duration - Duration of activity
  - `satisfaction`: Score(1-10) - Satisfaction with activity
- `washedHair`: Boolean - Whether hair was washed today
- `washedBody`: Boolean - Whether body was washed today

### WorkData
- `workDay`: Boolean - Whether it was a work day
- `workHours`: Duration - Hours worked
- `meetingTime`: Duration - Time spent in meetings
- `focusTime`: Duration - Time spent in focused work
- `workSatisfaction`: Score(1-10) - Satisfaction with work accomplished

### EnvironmentalData
- `weather`: WeatherMetrics - Data from weather API
  - `temperature`: Number - Average temperature
  - `humidity`: Number - Average humidity
  - `airPressure`: Number - Barometric pressure
  - `sunlightHours`: Number - Hours of daylight
  - `uvIndex`: Number - UV index for the day
  - `precipitation`: Number - Rainfall/snowfall amount
- `airQuality`: AirQualityMetrics - Data from air quality API
  - `aqi`: Number - Air Quality Index
  - `pm25`: Number - Fine particulate matter level
  - `pm10`: Number - Coarse particulate matter level
  - `o3`: Number - Ozone level
  - `no2`: Number - Nitrogen dioxide level
- `allergyIndex`: Number - Pollen/allergy index from API
- `moonPhase`: String - Current moon phase

## Implementation Phases

### Phase -1: LLM-Based Conversational Tracking
- **Daily Check-in Structure**:
  - **Morning Check-in**:
    - "How rested do you feel today? (1-10)"
    - "How happy do you feel this morning? (1-10)"
    - "How anxious do you feel this morning? (1-10)"
    - "What time did you go to bed and wake up?"
    - "Any supplements taken this morning?"

  - **Evening Check-in**:
    - "How happy do you feel tonight? (1-10)"
    - "How anxious do you feel tonight? (1-10)"
    - "How productive were you today? (1-10)"
    - "Did you work today or was it a free day?"
    - "What did you eat today? Any fast food or gluten?"
    - "Did you exercise today? Nature time? Step count?"
    - "Any supplements taken today? Creatine?"
    - "How much time did you spend on screens today?"
    - "How was your relationship with your partner today?"
    - "Any self-care activities today? Did you wash your hair or body today?"
    - "How much pressure do you feel from life admin? (1-10)"

- **LLM Prompt Template**:
  ```
  You are a life tracking assistant. Your job is to collect daily data points about my wellbeing,
  activities, and habits, then format this data in a structured JSON format for later analysis.

  IMPORTANT HISTORICAL EVENTS TO TRACK:
  - Last caffeine consumption: May 3, 2025 (calculate days since)
  - Started dutasteride: May 9, 2025 (track continued usage)
  - Last washed hair: May 8, 2025 (calculate days since)

  Ask me the appropriate morning or evening questions based on the current time.

  MORNING CHECK-IN QUESTIONS:
  - How rested do you feel today? (1-10)
  - How happy do you feel this morning? (1-10)
  - How anxious do you feel this morning? (1-10)
  - What time did you go to bed and wake up?
  - Any supplements taken this morning? (Ask specifically about dutasteride)

  EVENING CHECK-IN QUESTIONS:
  - How happy do you feel tonight? (1-10)
  - How anxious do you feel tonight? (1-10)
  - How productive were you today? (1-10)
  - Did you work today or was it a free day?
  - What did you eat today? Any fast food or gluten?
  - Did you exercise today? Nature time? Step count?
  - Any supplements taken today? (Ask specifically about dutasteride if not mentioned in morning)
  - Did you consume any caffeine today? (Tracking days since May 3)
  - How much time did you spend on screens today?
  - How was your relationship with your partner today?
  - Any self-care activities today? Did you wash your hair or body today?
  - How much pressure do you feel from life admin? (1-10)

  After I respond to all questions, create or update an artifact named "life-tracking-data" with the collected JSON data.
  If updating an existing artifact, append the new data as a new JSON object in an array.

  The JSON object should follow this structure:
  {
    "date": "YYYY-MM-DD",
    "timeOfDay": "morning|evening",
    "subjective": {
      "restedness": 1-10,
      "happiness": 1-10,
      "anxiety": 1-10,
      "productivity": 1-10,
      "lifeAdminPressure": 1-10
    },
    "sleep": {
      "bedTime": "HH:MM",
      "wakeTime": "HH:MM"
    },
    "nutrition": {
      "caffeineConsumed": boolean,
      "daysSinceCaffeine": number (calculated from May 3, 2025)
    },
    "supplements": {
      "dutasterideTaken": boolean,
      "daysSinceStarted": number (calculated from May 9, 2025)
    },
    "selfCare": {
      "washedHair": boolean,
      "washedBody": boolean
    },
    // other relevant fields based on answers...
    "derivedMetrics": {
      "daysSinceEvents": {
        "caffeine": number,
        "startedDutasteride": number,
        "washedHair": number,
        "washedBody": number
        // other events as tracked
      }
    }
  }

  Make the JSON valid but concise. Include only fields that have values.
  Store this data incrementally in the artifact, so we can build up a complete daily record.
  Remind me if I miss a check-in.
  ```

- Manual data collection through structured LLM conversations
- JSON data storage for easy migration to later phases
- Refinement of metrics and tracking categories
- Identification of key correlations and patterns
- Development of tracking habits before system implementation

### Phase 1: Mobile-First Basic System
- Go backend with HTMX frontend optimized for mobile
- Progressive Web App (PWA) for home screen installation
- Quick-tap interfaces for common inputs
- **Dynamic quick-action system** for frequently used inputs
- PostgreSQL data storage
- Basic visualization of recent trends
- Full Oura Ring API integration
- **Robust offline capabilities**:
  - Complete offline functionality for all data entry
  - Background synchronization when connection restored
  - Conflict resolution for offline changes
  - Local storage with encryption

### Phase 2: Integration & Automation
- Expanded API integrations (where available)
- Automated tracking for objective metrics
- Smart defaults and pattern recognition
- Notification system with actionable inputs
- Reminder system for daily inputs

### Phase 3: Analysis & Advanced Features
- Statistical analysis for identifying correlations
- Comprehensive dashboard across devices
- Advanced pattern recognition
- Recommendation engine
- Data export/sharing capabilities
- **Smart Recommendation System**:
  - Daily actionable insights based on sleep patterns
  - Specific "Do" and "Don't" recommendations for better sleep
  - Example: "Avoid AC tonight, wear compression socks, avoid eating within 3 hours of bedtime"
  - Personalized based on historical effectiveness of interventions
  - Presented as simple, actionable evening checklist

## Statistical Methodology

### Data Analysis Approach

#### Handling Missing Data
- **Interpolation Methods**:
  - Linear interpolation for continuous metrics (like sleep duration)
  - Last observation carried forward for categorical data
  - Multiple imputation for critical variables when appropriate
- **Missing Data Visualization**: Clear indication of interpolated/missing data in visualizations
- **Tracking Consistency Metrics**: Track and visualize data collection completeness

#### Correlation Analysis
- **Lag Analysis**: Examine effects across different time frames (same-day, next-day, multi-day)
- **Rolling Windows**: Use 3-day, 7-day, and 30-day windows to identify longer-term patterns
- **Multivariate Analysis**: Control for confounding variables when establishing correlations
  - Example: Sleep quality vs. mood, controlling for exercise and nutrition
- **Causality Testing**: Granger causality tests to suggest potential causal relationships

#### Confounding Variable Management
- **Factor Analysis**: Identify clusters of related variables
- **Stratification**: Analyze correlations within subsets of similar days
- **Control Variables**: Maintain awareness of external factors (weather, seasons, work demands)
- **Natural Experiments**: Leverage unplanned disruptions (travel, illness) to isolate variables

#### Insight Generation
- **Hypothesis Testing**: Formal testing of correlations with statistical significance
- **Anomaly Detection**: Identify outlier days and examine contributing factors
- **Pattern Recognition**: Machine learning techniques to identify complex patterns
  - Decision trees for identifying key decision points
  - Clustering for identifying similar "types of days"
- **Personalized Baselines**: Establish individual baselines rather than population norms

### Visualization Strategies
- **Correlation Heatmaps**: Visualize strength of relationships between variables
- **Time Series Analysis**: Track trends, seasonality, and cyclical patterns
- **Comparative Analysis**: Side-by-side comparison of different time periods
- **Causal Diagrams**: Visualize potential causal relationships between variables

## Technical Architecture

### Backend
- **Language**: Go
- **Web Framework**: [Appropriate Go web framework]
- **Database**:
  - PostgreSQL for structured data
  - Event store for all tracking events
  - **Architecture Note**: Implement event sourcing patterns where raw events are the source of truth, with materialized views for quick access
- **API**: RESTful endpoints + GraphQL (optional)
- **Authentication**: OAuth2 / JWT
- **Deployment**: render.com infrastructure

### Frontend
- **Primary Technology**: HTMX + Alpine.js
- **CSS Framework**: TailwindCSS
- **Mobile Optimization**: PWA with offline capabilities
  - Service workers for offline data collection
  - IndexedDB for local data storage
  - Background sync API for data transmission when online
  - Conflict resolution for simultaneous online/offline changes
- **Data Visualization**: Chart.js or D3.js
- **Quick Actions Engine**: Client/server system for dynamic action generation based on usage patterns

## Potential Data Integration Sources

The system can leverage numerous automated data sources to minimize manual tracking while maximizing insights. Below is a comprehensive list of potential integrations, organized by domain.

### Physical Activity & Health
1. **Oura Ring API** (Primary)
   - Sleep patterns, readiness scores, activity metrics
   - HRV, body temperature, and respiratory rate
   - Provides core sleep and recovery metrics

2. **Google Fit / Apple Health**
   - Steps, workouts, active minutes, heart rate
   - Aggregates data from multiple fitness apps
   - Provides holistic view of daily movement

3. **Additional Wearables**
   - Garmin, Fitbit, Whoop APIs
   - More detailed exercise and recovery metrics
   - Different emphasis on metrics compared to Oura

4. **Smart Scales**
   - Withings, Renpho, Eufy
   - Weight, body composition, body fat percentage
   - Long-term physical changes tracking

### Digital Behavior
5. **Android/iOS Screen Time**
   - Digital Wellbeing exports (Android)
   - Screen Time API (iOS)
   - App usage patterns and categories

6. **Browser Extensions**
   - RescueTime, Webtime Tracker, ManicTime
   - Detailed computer usage analytics
   - Productivity vs. distraction classification

7. **Email & Calendar**
   - Google Workspace / Microsoft Graph APIs
   - Meeting time, email volume, response times
   - Work patterns and potential stressors

8. **Social Media Usage**
   - Platform APIs where available
   - Usage patterns, posting frequency
   - Potential impact on mood and anxiety

9. **Gaming Platforms**
   - Steam API for game time tracking
   - PlayStation/Xbox activity if available
   - Game types, duration, timing patterns
   - Social vs. solo gaming habits

10. **Photo Library Analysis**
    - Google Photos API + Multimodal LLM analysis
    - Passive tracking of activities and social interactions
    - Location diversity and travel patterns
    - Nature exposure and outdoor activities
    - Facial expression analysis for mood indicators

### Location & Environment
9. **Location History**
   - Google Timeline data
   - Time at home, work, nature, social venues
   - Travel and commute patterns

10. **Weather APIs**
    - OpenWeatherMap, WeatherAPI, etc.
    - Sunlight exposure, temperature, pressure
    - Environmental factors affecting mood

11. **Air Quality Monitoring**
    - PurpleAir, AirNow, BreezoMeter APIs
    - Indoor and outdoor air quality metrics
    - Potential impacts on energy and respiratory health

12. **Home Environment**
    - Smart home sensors (temperature, humidity)
    - Light exposure patterns
    - Sleep environment quality indicators

13. **Allergy Index**
    - Pollen.com, AccuWeather, Weather.com Allergy APIs
    - Daily pollen counts and allergy forecast
    - Track correlation between allergies and sleep, energy, and productivity

### Consumption & Lifestyle
14. **Financial Data**
    - Bank/credit card export features
    - Food delivery, groceries, entertainment spending
    - Consumption patterns and potential stressors

15. **Media Consumption**
    - Spotify, Netflix, YouTube history APIs
    - Content types, duration, timing
    - Entertainment and information intake patterns

16. **Smart Home Devices**
    - Smart fridges, lighting systems, thermostats
    - Home activity patterns
    - Potential insights on daily routines

17. **Voice Assistant History**
    - Amazon Alexa or Google Home activity logs
    - Home automation patterns
    - Voice search trends and interests

### Importance and Implementation Priority

These integrations are valuable for different reasons:

- **Reducing Tracking Burden**: Automating data collection increases adherence
- **Objective Measurements**: Less susceptible to subjective biases than self-reporting
- **Continuous Monitoring**: Captures data even when user forgets to log
- **Rich Contextual Data**: Provides environmental and behavioral context for wellbeing patterns
- **Cross-Domain Insights**: Enables discovery of non-obvious correlations

**Implementation Priority**:
1. Core physical metrics (Oura, Phone Screen Time)
2. Digital behavior (Browser extensions, App usage)
3. Location and environmental factors (Weather, Air Quality, Allergy)
4. Consumption and lifestyle patterns

Each integration should include:
- Data normalization to fit our schema
- Privacy controls to limit sensitive data collection
- Backup mechanisms for when automated tracking fails

## Oura API Integration Specifications

### Authentication
- OAuth2 flow for authorization
- Token storage with refresh capabilities
- User permission for accessing:
  - Sleep data
  - Activity data
  - Readiness data

### Data Retrieval
- Daily synchronization job for latest data
- Endpoints to use:
  - `/v2/usercollection/sleep`
  - `/v2/usercollection/daily_activity`
  - `/v2/usercollection/daily_readiness`
- Key metrics to extract:
  - Sleep scores and stages
  - HRV and resting heart rate
  - Body temperature
  - Activity levels
  - Readiness score

### Processing
- Calculate sleep consistency metrics
- Normalize data for consistent tracking
- Generate derived metrics (e.g., consistency scores)

## User Experience Specifications

### Mobile Input Principles
- **One-Thumb Operation**: Design for single-handed use
- **Tap Over Type**: Prefer tap interactions over keyboard input
- **Smart Defaults**: Pre-fill with likely values based on patterns
- **Progressive Disclosure**: Hide complexity, reveal as needed
- **Batch Inputs**: Group related inputs to reduce context switching
- **Micro-Interactions**: Quick inputs throughout day vs. one long form
- **Dynamic Quick Actions**: Automatically generate buttons for recently/frequently used inputs

### Key Interfaces
- **Morning Check-In**:
  - Quick subjective scores (restedness, happiness, anxiety)
  - Confirm/adjust sleep timings from Oura
  - Previous day recap for missing data
  - Dynamic quick-action buttons for common morning activities/supplements

- **Evening Check-In**:
  - Meal tracking interface with quick categories
  - Supplement tracking with favorites system
  - Relationship and interaction tracking
  - Self-care activity logging
  - Work productivity assessment
  - Dynamic quick-action buttons for common evening activities

- **Ad-Hoc Tracking**:
  - Quick add for activities, meals, supplements
  - Photo input option for meals
  - Voice note capability for qualitative data
  - Recent/frequent items as tappable buttons (top 10-20)
  - Contextual suggestions based on time of day and patterns

- **Action Logging System**:
  - Reverse-chronological timeline of all logged actions
  - Daily summary of everything logged
  - Missing data alerts for categories not yet logged
  - Easily accessible with simple swipe gesture from main interface
  - Option to edit or delete recently logged items

## Next Steps

1. **Immediate Data Collection Setup**
   - Begin collecting data through the LLM-based system
   - Develop simple logging interface for the current phase
   - Analyze initial data to refine metrics and tracking categories

2. **Oura Integration**
   - Set up Oura API authentication
   - Develop basic data import functionality
   - Create data normalization processes

3. **PWA Development**
   - Set up repository structure
   - Implement core domain models in Go
   - Create database schema on render.com
   - Develop mobile-first UI with HTMX
   - Build quick-actions engine

4. **Analysis Capabilities**
   - Implement basic correlation analysis
   - Develop visualization dashboard
   - Create daily/weekly report generation

## File Structure

```
/cmd
  /api        # API server entry point
  /migrations # Database migrations
/internal
  /domain
    /model      # Domain model definitions
    /event      # Event definitions and handlers
    /aggregate  # Aggregate models built from events
  /service    # Business logic implementation
  /repository
    /event      # Event storage interfaces
    /view       # Materialized view storage
  /api        # API handlers
  /integration # Third-party integrations
    /oura     # Oura Ring API integration
    /weather  # Weather API integration
    /allergy  # Allergy API integration
  /analysis   # Data analysis capabilities
/web
  /templates  # HTMX templates
  /static     # Static assets
  /js         # Client-side JavaScript
/scripts      # Utility scripts
/docs         # Documentation
```

### Event Architecture Details

The system will use an event-sourcing approach where:

1. **Events as Source of Truth**:
   - All user actions and integrations generate immutable events
   - Example events: `SupplementTaken`, `MealConsumed`, `ExerciseCompleted`, `SleepRecorded`
   - Each event includes timestamp, type, and payload with relevant details
   - Events are never updated or deleted, only appended

2. **Event Storage**:
   - Events stored chronologically in an append-only log
   - Optimized for write performance
   - Designed for easy replay and regeneration of state

3. **Materialized Views**:
   - Daily records generated by aggregating relevant events
   - Cached for performance but can be regenerated from events
   - Updated in real-time as new events are received

4. **Benefits of this Approach**:
   - Complete history and audit trail of all data
   - Flexible data model that can evolve over time
   - Easier to generate different aggregation timeframes (daily, weekly, monthly)
   - Natural support for "days since" metrics by querying most recent event of a type
   - Simplifies correction of historical data

5. **Implementation Considerations**:
   - Use PostgreSQL with JSONB for event storage
   - Implement event handlers to update materialized views
   - Design efficient indexing for event querying
   - Implement trigger-based refreshes for materialized views
   - OR Add inline application refreshes after data modifications

## Limitations and Considerations

### Data Collection Challenges
- **Tracking Fatigue Risk**: The comprehensive nature of metrics may lead to inconsistent tracking
- **Solution**: Implement a tiered approach with "core" vs "optional" metrics
- **Consistency Challenges**: Multiple input methods could create data fragmentation
- **Solution**: Define strict data schemas and normalization processes

### Privacy and Security
- **Sensitive Data**: System contains personal health and relationship information
- **Data Protection**: Implement appropriate encryption and access controls
- **Data Ownership**: All data remains user-owned with complete export capabilities
- **Backup Strategy**: Regular automated backups with versioning

### Statistical Analysis Considerations
- **Confounding Variables**: Many life factors correlate without causation
- **Missing Data Handling**: Statistical methods to account for tracking gaps
- **Insight Validation**: Mechanisms to validate and test correlation hypotheses

### Technical Considerations
- **Offline Capability**: PWA must function without internet connection
- **Data Synchronization**: Conflict resolution for offline data
- **Progressive Enhancement**: Core functionality works without advanced features
- **Testing Strategy**: Automated tests for data integrity and UI functionality

### User Experience Challenges
- **Consistency vs. Convenience**: Balance between detailed data and tracking ease
- **Notification Management**: Avoid alert fatigue while maintaining tracking consistency
- **Learning Curve**: Interface complexity increases with feature richness
- **Solution**: Implement progressive disclosure of advanced features

## Theories and Hypotheses to Test

### Sleep Quality Theories
1. **Hot Shower Effect Theory**
   - **Hypothesis**: Taking a hot shower before bed improves sleep quality by reducing time to fall asleep and increasing deep sleep
   - **Metrics to Compare**:
     - Sleep score with/without hot shower
     - Deep sleep percentage with/without hot shower
     - Restedness score the next morning
   - **Control Variables**: Room temperature, screen time before bed

2. **Hungry at Bedtime Theory**
   - **Hypothesis**: Going to bed slightly hungry improves sleep quality and next-day energy levels
   - **Metrics to Compare**:
     - Sleep score when hungry vs. full
     - HRV during sleep with different meal timings
     - Morning restedness based on hours since last meal
   - **Control Variables**: Meal content, caffeine intake, sleep timing

3. **Grey Noise & Noise Pollution Theory**
   - **Hypothesis**: Grey noise eliminates the negative effects of environmental noise pollution on sleep
   - **Metrics to Compare**:
     - Awake time during sleep with/without grey noise in noisy conditions
     - Sleep continuity score with/without grey noise
     - Subjective restedness based on noise conditions and interventions
   - **Control Variables**: Overall sleep duration, earplugs usage, partner presence

4. **Screen Exposure Timing Theory**
   - **Hypothesis**: Screen exposure within 60 minutes of sleep significantly reduces deep sleep percentage
   - **Metrics to Compare**:
     - Deep sleep percentage with different screen cutoff times
     - Sleep latency (time to fall asleep) with different screen exposure patterns
   - **Control Variables**: Blue light filtering, content type (work vs. entertainment)

5. **Compression Socks Theory**
   - **Hypothesis**: Wearing compression socks during sleep improves circulation and increases deep sleep
   - **Metrics to Compare**:
     - Deep sleep percentage with/without socks
     - HRV patterns throughout the night
     - Morning restedness and leg comfort
   - **Control Variables**: Room temperature, blanket weight, mattress type

6. **Clothing Restriction Theory**
   - **Hypothesis**: Sleeping without a shirt improves sleep quality by allowing better temperature regulation
   - **Metrics to Compare**:
     - Sleep quality with/without shirt
     - Temperature regulation (via Oura body temperature)
     - Number of awakenings during night
   - **Control Variables**: Room temperature, blanket type, AC usage

7. **Fresh Bed Sheets Theory**
   - **Hypothesis**: Sleeping on fresh, clean bed sheets improves sleep quality and restedness
   - **Mechanism**: May reduce exposure to allergens, dust mites, bed bugs, or accumulated dead skin cells that could trigger allergic responses or mild immune system activation during sleep
   - **Metrics to Compare**:
     - Sleep quality in days after changing sheets vs. days before
     - Subjective restedness scores when sleeping on fresh sheets
     - Allergy symptoms and breathing quality
     - Deep sleep percentage and sleep fragmentation
   - **Control Variables**: Other sleep conditions, season, detergent used
   - **Related Hypotheses**: Regular mattress cleaning and pillow replacement may show similar benefits through allergen reduction

8. **Supplement Timing Theory**
   - **Hypothesis**: Taking supplements consistently at the same time daily increases their effectiveness
   - **Metrics to Compare**:
     - Wellbeing scores with consistent vs. inconsistent supplement timing
     - Long-term trends in measured metrics based on supplement consistency
   - **Control Variables**: Supplement types, dosages, other lifestyle factors

8. **Self-Gratification Effect Theory**
   - **Hypothesis**: Self-gratification before bed improves sleep onset and quality through relaxation
   - **Metrics to Compare**:
     - Time to fall asleep with/without self-gratification
     - Sleep continuity throughout the night
     - Morning restedness scores
   - **Control Variables**: Content usage during self-gratification, time before sleep, other relaxation techniques

### Productivity & Energy Theories
9. **Social Media Impact Theory**
   - **Hypothesis**: Morning social media usage reduces productivity more than evening usage
   - **Metrics to Compare**:
     - Productivity scores on days with morning vs. evening social media usage
     - Focus duration and work satisfaction based on social media timing
   - **Control Variables**: Total screen time, work type, sleep quality

10. **Caffeine Elimination Theory**
    - **Hypothesis**: Complete elimination of caffeine improves sleep scores and reduces anxiety after a 14-day adaptation period
    - **Metrics to Compare**:
      - Sleep quality, anxiety, and HRV trends following caffeine elimination
      - Return to baseline after reintroduction
    - **Control Variables**: Exercise, hydration, stress levels

11. **Nature Time Dosage Theory**
    - **Hypothesis**: At least 20 minutes of nature time produces measurable improvements in wellbeing scores
    - **Metrics to Compare**:
      - Wellbeing and anxiety scores based on nature exposure duration
      - Cumulative effect of consecutive days with nature time
    - **Control Variables**: Weather conditions, exercise during nature time

12. **Environmental Allergy Impact Theory**
    - **Hypothesis**: High allergy index days reduce productivity and sleep quality
    - **Metrics to Compare**:
      - Productivity and sleep scores correlated with allergy index
      - Effectiveness of interventions (air purifiers, medications) during high allergy periods
    - **Control Variables**: Indoor vs. outdoor time, medication use

## Original Tracking Requirements

The following is the original set of tracking requirements that informed this specification:

```
Together I want us to build something that will let me tell it how I feel, what I take, how I've slept, with a general formula of trying to figure out: what i took yesterday, how i felt, how much i exercised, how much i worked, how much stress = how much sleep + how good i feel the next day

I want to track things like, in order to import data somewhere useful so we can analyze it:
* All Oura Ring data on my sleep, activities, readiness, step count
* How do I feel on wake up (rested score, happiness score, anxiety score)
* How did I feel today (Rested score, happiness score, anxiety score)
* How productive I felt at work today
* Whether I worked, or am off (weekends count as off)
* Auto track, how long it's been:
   * since I worked
   * since free day
   * since at least 1w holiday
   * since fast food
   * since gluten, maybe other potential toxins
   * since called family
   * since exercise, since at least 4k steps
   * since XX supplement
   * since caffeine
   * since started taking dutasteride, minoxidil
   * since I was last in nature
   * since interacted with animals (cat, dog, etc)
   * since I worked out (eg push ups count)
   * since partner was next to me
   * since slept with partner
   * since (big) fight with partner
   * since washed hair
   * since washed body
   * (some of this stuff means I also need to track daily, like whether i was in nature or not! make sure to understand that as well)
* How consistent was the time in going to bed today, and 7d, 30d average
* How consistent was wake up time today, and 7d, 30d average
* How many hours did I eat before sleep?
* What type of food/drinks i ate (discrete categories, eg fast food, salad, mexican, sushi, etc...)
   * Whether there was gluten (or other potential toxins) (discrete categories)
   * Maybe instead listing ingredients also? can be auto determined eg if burger
   * for breakfast (if had any), lunch, dinner
   * how much satiety (1-10), 10 = overate
* How much creatine I took, when
* What supplements I took, when
* How much screen time (work comp, personal comp, and android: instagram, phone, youtube)
* How much time spent outside
* How much pressure do I feel from life admin
* How today's relationship with partner
```
