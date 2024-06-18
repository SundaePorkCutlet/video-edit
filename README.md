# ffmpeg를 이용한 동영상 편집 서버
동영상 업로드, 트림(Trim), 이어 붙이기(Concat), 정보 조회 기능을 제공합니다.  
  
# Result Code
**200** : 요청 성공  
**400** : 요청 오류  
**500** : 서버 오류  

**1001** : Trim 성공  
**1002** : Concat 성공  
**1003** : Trim 성공 Concat 실패  
**1004** : Concat 성공 Trim 실패  
**1005** : Trim Concat 성공  

**6001** : 허용되지 않는 확장자  
**6002** : Trim 실패  
**6003** : Concat 실패    

**example**  
  
<pre><code>{
    "code": 6002
}</code></pre>


# Requirements
- ***make***: [Installation Guide](https://www.gnu.org/software/make/)
- ***go 1.22***: [Installation Guide](https://golang.org/doc/install)
- ***MariaDB 11.2.2***: [Installation Guide](https://mariadb.org/download/)



# Script
### `make run`  
9000번 포트로 프로젝트를 실행시킵니다.  
localhost:9000으로 기본 설정이 되어있습니다.  
app_config.yaml 파일에서 포트와 주소를 변경할 수 있습니다.  



# API Example
### POST
### `/upload`  
동영상을 업로드합니다.  
Multipart form 형식으로 업로드합니다.   
app_config.yaml에 설정된 위치에 저장됩니다.    

**postman example**
  
<img width="834" alt="image" src="https://github.com/SundaePorkCutlet/video-edit/assets/87690981/9ae267ef-1423-4a18-b6d2-a8f09df11319">
  
### POST  
### `/modify`
3가지 조건에 따라 파라미터가 변경됩니다.  


### 1. trim만 요청하는 경우
<pre><code>{
    "isTrimed":true,
    "trimVideoList":[
        {
        "videoId":"2d05c690-ee19-49d2-b256-85a1641cd24b",
        "startTime":2,
        "endTime":5
        },
        {
        "videoId":"210a0287-2d83-4b80-afa5-44cdd709650d",
        "startTime":1,
        "endTime":70
        }
    ]
}</code></pre>

isTrimed : true로 설정
trimVideoList : 여러 동영상 트림 요청 가능 (배열 형식)  
videoId : 요청할 동영상 고유 id  
startTime : 시작 시간 (초 단위)  
endTime : 종료 시간 (초 단위)  

DB의 trim_history 테이블에 생성된 동영상 UUID와 원본 동영상 UUID, startTime, endTime이 저장됩니다.  

### 2. concat만 요청하는 경우
<pre><code>{
    "isConcated":true,
    "concatVideoList":[
        "0c649117-6ae9-4d22-b033-80c646191fa5","c801f057-d346-4889-8d5d-68878423f03e"
    ]
}</code></pre>  

  isConcated : true로 설정  
  concatVideoList : 요청할 동영상 고유 id 리스트  

app_config.yaml에 설정된 위치에 원본 동영상 경로 리스트가 담긴 txt 파일이 생성됩니다.  
DB의 concat_history 테이블에 생성된 동영상 UUID와 인코딩된 동영상 UUID 리스트 txt 파일 경로가 저장됩니다.  

### 3. trim한 동영상들을 concat 요청하는 경우  
<pre><code>{
    "isTrimed":true,
    "trimVideoList":[
         {
        "videoId":"7d04454e-5fd9-48b8-a8bd-ebdabd7e3362",
        "startTime":2,
        "endTime":5
        },
        {
        "videoId":"83158ed9-5ff4-46b2-aed9-63064196b505",
        "startTime":1,
        "endTime":7
        }
    ],
    "isConcated":true
}</code></pre>  

isConcated가 true로 추가 설정됩니다.  
DB의 trim_history 테이블에는 생성된 동영상 UUID와 원본 동영상 UUID, startTime, endTime이 저장됩니다.  
DB의 encode_history 테이블에는 인코딩된 동영상 UUID와 Trim된 원본 동영상 UUID가 저장됩니다.  
DB의 concat_history 테이블에는 생성된 동영상 UUID와 인코딩된 동영상 UUID 리스트 txt 파일 경로가 저장됩니다.  
  
### GET
### `/video`  
동영상 모든 정보를 제공합니다  




#### result example
<pre><code>[
    {
        "video": {
            "id": "105251f7-7516-42fc-80fd-86f0f0ba5368",
            "path": "/Users/hongjunho/Downloads/workspace/stockfolio-test/file/videos/105251f7-7516-42fc-80fd-86f0f0ba5368.mp4",
            "videoName": "trim_KakaoTalk_Video_2024-06-14-20-28-32.mp4",
            "extension": "mp4",
            "uploadTime": "",
            "isTrimed": true,
            "trimTime": "20240617212941",
            "isConcated": false,
            "concatTime": "",
            "isEncoded": false,
            "encodeTime": ""
        },
        "trimInfo": {
            "videoId": "7d04454e-5fd9-48b8-a8bd-ebdabd7e3362",
            "videoPath": "/Users/hongjunho/Downloads/workspace/stockfolio-test/file/videos/7d04454e-5fd9-48b8-a8bd-ebdabd7e3362.mp4",
            "startTime": 2,
            "endTime": 5
        },
        "concatInfoPath": "",
        "encodeInfoPath": ""
    },
]</code></pre>  

### GET
### `/download/:uuid`  
동영상 다운로드 제공 안내

동영상 파일을 다운로드할 때의 파일 이름 지정 방식은 다음과 같습니다:

1. **원본 동영상**: 원본 동영상의 이름으로 제공됩니다.
2. **Trim된 파일**: `tr-{원본동영상이름}.{원본확장자}` 형식으로 제공됩니다.
3. **인코딩된 파일**:
   - **Concat 파일**: 인코딩된 파일은 `cv-{원본동영상이름}.mp4` 형식으로 제공됩니다.
   - **Trim된 파일을 인코딩한 경우**: `cv-tr-{원본동영상이름}.mp4` 형식으로 제공됩니다.
4. **Concat된 파일**: `cc-{concat 요청시간}.mp4` 형식으로 제공됩니다.
5. **Trim과 Concat을 모두 요청한 파일**: Concat 파일 형식(`cc-{concat 요청시간}.mp4`)으로 제공됩니다.
  

### example  
#### Trim과 인코딩된 동영상 다운로드 요청 시

- **파일 이름**: `cv-tr-{원본동영상이름}.mp4`

  
<img width="648" alt="image" src="https://github.com/SundaePorkCutlet/video-edit/assets/87690981/1eac4062-6325-4db7-bb3c-25c211823a58">    

<img width="330" alt="image" src="https://github.com/SundaePorkCutlet/video-edit/assets/87690981/31c41461-f086-4628-868a-1d5cbdd599c9">    

#### Concat된 동영상 다운로드 요청 시

- **파일 이름**: `cc-{concat 요청시간}.mp4`
  
<img width="338" alt="image" src="https://github.com/SundaePorkCutlet/video-edit/assets/87690981/de3e8927-84e0-4321-a79d-e88a4af50dc3">  




## ISSUE
동영상을 Concat 하기 위해 해상도, 확장자, 코덱을 통일해야 했기 때문에 Concat하기 전 모든 동영상을 인코딩합니다.  
이로 인해 Trim한 동영상들을 Concat 요청 시 동영상 수가 급격히 증가하는 문제가 있습니다.  
ex) 3개의 동영상을 trim시킨다음에 concat을 요청함  
1. 세 개의 동영상을 Trim한 뒤 저장 -> 동영상 3개 증가, trim_history에 Trim된 동영상 UUID와 원본 동영상 UUID 저장    
2. Trim된 세 개의 동영상을 Concat하기 위해 인코딩 -> 동영상 3개 증가, encode_history에 인코딩된 동영상 UUID와 Trim된 동영상 UUID 저장    
3. 인코딩된 세 개의 동영상을 Concat -> 동영상 1개 증가, Concat txt 파일 생성 및 인코딩된 동영상 리스트 생성, concat_history에 Concat된 동영상 UUID와 txt 파일 경로 저장   
