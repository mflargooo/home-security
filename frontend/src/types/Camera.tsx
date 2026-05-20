export type CameraStatus =
    | 'active'
    | 'reconnecting'
    | 'offline'
    | 'disabled'

export type Camera = {
	id:           string,            
	name:         string,            
	rtsp_url:     string,            
	status:       CameraStatus,      
	metadata?:    Record<string, string>, 
	created_at:   string,         
	updated_at:   string,         
	last_seen_at: string | null,        
}